/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	billingRepository "github.com/coze-dev/coze-studio/backend/domain/billing/repository"
)

// ============================================================
// 提现管理服务 - 基于企业级支付规范
// ============================================================

// WithdrawalService 提现服务
type WithdrawalService struct {
	db                *gorm.DB
	withdrawalRepo    billingRepository.WithdrawalRequestRepository
	limitRepo         billingRepository.WithdrawalLimitRepository
	recordRepo        billingRepository.WithdrawalRecordRepository
	reconciliationRepo billingRepository.WithdrawalReconciliationRepository
	billingAccountRepo billingRepository.BillingAccountRepository
	logger            *logrus.Logger

	// 加密密钥（用于银行卡号加密）
	encryptionKey     []byte
}

// NewWithdrawalService 创建提现服务实例
func NewWithdrawalService(
	db *gorm.DB,
	withdrawalRepo billingRepository.WithdrawalRequestRepository,
	limitRepo billingRepository.WithdrawalLimitRepository,
	recordRepo billingRepository.WithdrawalRecordRepository,
	reconciliationRepo billingRepository.WithdrawalReconciliationRepository,
	billingAccountRepo billingRepository.BillingAccountRepository,
	encryptionKey string,
	logger *logrus.Logger,
) *WithdrawalService {
	return &WithdrawalService{
		db:                db,
		withdrawalRepo:    withdrawalRepo,
		limitRepo:         limitRepo,
		recordRepo:        recordRepo,
		reconciliationRepo: reconciliationRepo,
		billingAccountRepo: billingAccountRepo,
		encryptionKey:     []byte(encryptionKey),
		logger:            logger,
	}
}

// ============================================================
// 请求和响应类型
// ============================================================

// CreateWithdrawalRequest 创建提现申请请求
type CreateWithdrawalRequest struct {
	TenantID         string  `json:"tenant_id" validate:"required"`
	Amount           float64 `json:"amount" validate:"required,gt=0"`
	WithdrawalMethod string  `json:"withdrawal_method" validate:"required"` // alipay, wechat, bank_transfer
	BankAccount      string  `json:"bank_account,omitempty"`
	BankName         *string `json:"bank_name,omitempty"`
	AccountName      *string `json:"account_name,omitempty"`
	AccountType      string  `json:"account_type" validate:"required"` // personal, company
	Notes            *string `json:"notes,omitempty"`
}

// WithdrawalRequestResponse 提现申请响应
type WithdrawalRequestResponse struct {
	WithdrawalID     uint64  `json:"withdrawal_id"`
	WithdrawalNumber string  `json:"withdrawal_number"`
	Amount           float64 `json:"amount"`
	Fee              float64 `json:"fee"`
	ActualAmount     float64 `json:"actual_amount"`
	Currency         string  `json:"currency"`
	Status           string  `json:"status"`
	EstimatedArrivalTime int64 `json:"estimated_arrival_time"`
	CreatedAt        int64   `json:"created_at"`
}

// ReviewWithdrawalRequest 审核提现申请请求
type ReviewWithdrawalRequest struct {
	WithdrawalID  uint64  `json:"withdrawal_id" validate:"required"`
	Approved      bool    `json:"approved" validate:"required"`
	ReviewComment *string `json:"review_comment,omitempty"`
	ReviewerID    uint64  `json:"reviewer_id" validate:"required"`
}

// WithdrawalStatisticsResponse 提现统计响应
type WithdrawalStatisticsResponse struct {
	TenantID        string  `json:"tenant_id"`
	Period          *Period `json:"period,omitempty"`
	TotalCount      int     `json:"total_count"`
	TotalAmount     float64 `json:"total_amount"`
	SuccessCount    int     `json:"success_count"`
	SuccessAmount   float64 `json:"success_amount"`
	PendingCount    int     `json:"pending_count"`
	PendingAmount   float64 `json:"pending_amount"`
	FailedCount     int     `json:"failed_count"`
	FailedAmount    float64 `json:"failed_amount"`
	TotalFee        float64 `json:"total_fee"`
	AverageAmount   float64 `json:"average_amount"`
	SuccessRate     float64 `json:"success_rate"`
}

// Period 时间周期
type Period struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// ============================================================
// 核心方法实现
// ============================================================

// CreateWithdrawal 创建提现申请
// 参数:
//   - ctx: 上下文
//   - req: 创建提现申请请求
// 返回: 提现申请响应或错误
func (s *WithdrawalService) CreateWithdrawal(
	ctx context.Context,
	req *CreateWithdrawalRequest,
) (*WithdrawalRequestResponse, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		RecordWithdrawalOperation("create", req.WithdrawalMethod, duration)
	}()

	s.logger.WithFields(logrus.Fields{
		"tenant_id":          req.TenantID,
		"amount":             req.Amount,
		"withdrawal_method":  req.WithdrawalMethod,
	}).Info("[WithdrawalService] Creating withdrawal request")

	// 1. 验证金额
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	// 2. 查询计费账户
	account, err := s.billingAccountRepo.GetByTenant(ctx, req.TenantID)
	if err != nil {
		s.logger.WithError(err).Errorf("[WithdrawalService] Failed to get billing account")
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	if account == nil {
		return nil, fmt.Errorf("no billing account found for tenant: %s", req.TenantID)
	}

	// 3. 检查余额
	if account.AvailableCredit < req.Amount {
		return nil, fmt.Errorf("insufficient balance: available=%.2f, requested=%.2f",
			account.AvailableCredit, req.Amount)
	}

	// 4. 获取限额配置
	limit, err := s.limitRepo.GetEffectiveLimit(ctx, req.WithdrawalMethod, req.Amount)
	if err != nil {
		s.logger.WithError(err).Errorf("[WithdrawalService] Failed to get withdrawal limit")
		return nil, fmt.Errorf("failed to get withdrawal limit: %w", err)
	}

	// 5. 验证限额
	if err := s.validateWithdrawalLimit(ctx, req.TenantID, req.Amount, limit); err != nil {
		return nil, err
	}

	// 6. 计算手续费
	fee := s.calculateFee(req.Amount, limit)

	// 7. 验证总金额（提现金额 + 手续费）
	totalAmount := req.Amount + fee
	if account.AvailableCredit < totalAmount {
		return nil, fmt.Errorf("insufficient balance for withdrawal including fee: available=%.2f, required=%.2f",
			account.AvailableCredit, totalAmount)
	}

	// 8. 生成提现单号
	withdrawalNumber, err := s.generateWithdrawalNumber(ctx, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate withdrawal number: %w", err)
	}

	// 9. 风控检查
	riskScore, riskReason := s.performRiskCheck(ctx, req, account)

	// 10. 加密银行账号
	encryptedAccount, err := s.encryptBankAccount(req.BankAccount)
	if err != nil {
		s.logger.WithError(err).Errorf("[WithdrawalService] Failed to encrypt bank account")
		return nil, fmt.Errorf("failed to encrypt bank account: %w", err)
	}

	// 11. 创建提现申请
	withdrawal := &entity.WithdrawalRequest{
		TenantID:         req.TenantID,
		BillingAccountID: account.ID,
		WithdrawalNumber: withdrawalNumber,
		Amount:           req.Amount,
		Fee:              fee,
		ActualAmount:     req.Amount - fee,
		Currency:         account.Currency,
		WithdrawalMethod: req.WithdrawalMethod,
		BankAccount:      encryptedAccount,
		BankName:         req.BankName,
		AccountName:      req.AccountName,
		AccountType:      req.AccountType,
		Provider:         s.getProvider(req.WithdrawalMethod),
		Status:           entity.WithdrawalStatusPending,
		RiskScore:        &riskScore,
		RiskReason:       &riskReason,
		Notes:            req.Notes,
	}

	// 自动审核逻辑
	if limit != nil && limit.AutoApproveLimit != nil && req.Amount <= *limit.AutoApproveLimit {
		withdrawal.Status = entity.WithdrawalStatusApproved
		withdrawal.AutoApproved = true
		withdrawal.ReviewedAt = &startTime
	}

	// 12. 事务处理
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 12.1 保存提现申请
		if err := s.withdrawalRepo.Create(ctx, withdrawal); err != nil {
			s.logger.WithError(err).Errorf("[WithdrawalService] Failed to save withdrawal request")
			return fmt.Errorf("failed to save withdrawal request: %w", err)
		}

		// 12.2 冻结余额
		newBalance := account.AvailableCredit - totalAmount
		if err := s.billingAccountRepo.UpdateBalance(ctx, req.TenantID, -totalAmount); err != nil {
			return fmt.Errorf("failed to freeze balance: %w", err)
		}

		// 12.3 创建流水记录
		record := &entity.WithdrawalRecord{
			TenantID:         req.TenantID,
			WithdrawalID:     withdrawal.ID,
			WithdrawalNumber: withdrawalNumber,
			TransactionType:  entity.WithdrawalTransactionTypeFreeze,
			Amount:           -totalAmount,
			BalanceBefore:    account.AvailableCredit,
			BalanceAfter:     newBalance,
			Currency:         account.Currency,
			Description:      fmt.Sprintf("提现冻结：提现单号=%s, 金额=%.2f, 手续费=%.2f", withdrawalNumber, req.Amount, fee),
			TransactionDate:  time.Now().Format("2006-01-02"),
		}
		if err := s.recordRepo.Create(ctx, record); err != nil {
			return fmt.Errorf("failed to create withdrawal record: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 13. 记录日志
	s.logger.WithFields(logrus.Fields{
		"tenant_id":           req.TenantID,
		"withdrawal_id":       withdrawal.ID,
		"withdrawal_number":   withdrawalNumber,
		"amount":              req.Amount,
		"fee":                 fee,
		"actual_amount":       withdrawal.ActualAmount,
		"status":              withdrawal.Status,
		"auto_approved":       withdrawal.AutoApproved,
		"risk_score":          riskScore,
		"creation_time_s":     time.Since(startTime).Seconds(),
	}).Info("[WithdrawalService] Withdrawal request created successfully")

	// 14. 更新Prometheus指标
	UpdateWithdrawalMetrics(req.TenantID, req.WithdrawalMethod, "created", req.Amount)

	// 15. 计算预计到账时间（T+1或T+0）
	estimatedArrivalTime := s.calculateEstimatedArrivalTime(withdrawal)

	return &WithdrawalRequestResponse{
		WithdrawalID:          withdrawal.ID,
		WithdrawalNumber:      withdrawalNumber,
		Amount:                req.Amount,
		Fee:                   fee,
		ActualAmount:          withdrawal.ActualAmount,
		Currency:              account.Currency,
		Status:                withdrawal.Status,
		EstimatedArrivalTime:  estimatedArrivalTime,
		CreatedAt:             withdrawal.CreatedAt.Unix(),
	}, nil
}

// ReviewWithdrawal 审核提现申请
// 参数:
//   - ctx: 上下文
//   - req: 审核请求
// 返回: 错误
func (s *WithdrawalService) ReviewWithdrawal(
	ctx context.Context,
	req *ReviewWithdrawalRequest,
) error {
	startTime := time.Now()

	s.logger.WithFields(logrus.Fields{
		"withdrawal_id": req.WithdrawalID,
		"approved":      req.Approved,
		"reviewer_id":   req.ReviewerID,
	}).Info("[WithdrawalService] Reviewing withdrawal request")

	// 1. 查询提现申请
	withdrawal, err := s.withdrawalRepo.GetByID(ctx, req.WithdrawalID)
	if err != nil {
		s.logger.WithError(err).Errorf("[WithdrawalService] Failed to get withdrawal request")
		return fmt.Errorf("failed to get withdrawal request: %w", err)
	}

	if withdrawal == nil {
		return fmt.Errorf("withdrawal request not found: %d", req.WithdrawalID)
	}

	// 2. 检查状态
	if withdrawal.Status != entity.WithdrawalStatusPending &&
	   withdrawal.Status != entity.WithdrawalStatusReviewing {
		return fmt.Errorf("withdrawal request cannot be reviewed in current status: %s", withdrawal.Status)
	}

	// 3. 事务处理
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 3.1 更新提现状态
		now := time.Now()
		withdrawal.ReviewedAt = &now
		withdrawal.ReviewedBy = &req.ReviewerID
		withdrawal.ReviewComment = req.ReviewComment

		if req.Approved {
			withdrawal.Status = entity.WithdrawalStatusApproved
		} else {
			withdrawal.Status = entity.WithdrawalStatusRejected

			// 3.2 解冻余额
			totalAmount := withdrawal.Amount + withdrawal.Fee
			if err := s.billingAccountRepo.UpdateBalance(ctx, withdrawal.TenantID, totalAmount); err != nil {
				return fmt.Errorf("failed to unfreeze balance: %w", err)
			}

			// 3.3 创建流水记录
			account, _ := s.billingAccountRepo.GetByTenant(ctx, withdrawal.TenantID)
			record := &entity.WithdrawalRecord{
				TenantID:         withdrawal.TenantID,
				WithdrawalID:     withdrawal.ID,
				WithdrawalNumber: withdrawal.WithdrawalNumber,
				TransactionType:  entity.WithdrawalTransactionTypeUnfreeze,
				Amount:           totalAmount,
				BalanceBefore:    account.AvailableCredit - totalAmount,
				BalanceAfter:     account.AvailableCredit,
				Currency:         withdrawal.Currency,
				Description:      fmt.Sprintf("提现拒绝，解冻金额：提现单号=%s", withdrawal.WithdrawalNumber),
				TransactionDate:  now.Format("2006-01-02"),
			}
			if err := s.recordRepo.Create(ctx, record); err != nil {
				return fmt.Errorf("failed to create withdrawal record: %w", err)
			}
		}

		// 3.4 保存更新
		if err := s.withdrawalRepo.Update(ctx, withdrawal); err != nil {
			return fmt.Errorf("failed to update withdrawal request: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 4. 记录日志
	s.logger.WithFields(logrus.Fields{
		"withdrawal_id":   req.WithdrawalID,
		"status":          withdrawal.Status,
		"review_time_s":   time.Since(startTime).Seconds(),
	}).Info("[WithdrawalService] Withdrawal request reviewed successfully")

	// 5. 更新Prometheus指标
	status := "approved"
	if !req.Approved {
		status = "rejected"
	}
	UpdateWithdrawalMetrics(withdrawal.TenantID, withdrawal.WithdrawalMethod, status, withdrawal.Amount)

	return nil
}

// ExecuteWithdrawal 执行提现
// 参数:
//   - ctx: 上下文
//   - withdrawalID: 提现ID
// 返回: 错误
func (s *WithdrawalService) ExecuteWithdrawal(
	ctx context.Context,
	withdrawalID uint64,
) error {
	startTime := time.Now()

	s.logger.WithField("withdrawal_id", withdrawalID).
		Info("[WithdrawalService] Executing withdrawal")

	// 1. 查询提现申请
	withdrawal, err := s.withdrawalRepo.GetByID(ctx, withdrawalID)
	if err != nil {
		return fmt.Errorf("failed to get withdrawal request: %w", err)
	}

	if withdrawal == nil {
		return fmt.Errorf("withdrawal request not found: %d", withdrawalID)
	}

	// 2. 检查状态
	if withdrawal.Status != entity.WithdrawalStatusApproved {
		return fmt.Errorf("withdrawal request is not approved: %s", withdrawal.Status)
	}

	// 3. 更新状态为处理中
	withdrawal.Status = entity.WithdrawalStatusProcessing
	now := time.Now()
	withdrawal.ProcessedAt = &now
	if err := s.withdrawalRepo.Update(ctx, withdrawal); err != nil {
		return fmt.Errorf("failed to update withdrawal status: %w", err)
	}

	// 4. 执行提现（调用第三方支付接口）
	transactionID, err := s.executeWithdrawalViaProvider(ctx, withdrawal)
	if err != nil {
		s.logger.WithError(err).Errorf("[WithdrawalService] Failed to execute withdrawal via provider")

		// 4.1 更新为失败状态
		withdrawal.Status = entity.WithdrawalStatusFailed
		reason := err.Error()
		withdrawal.FailedReason = &reason
		if updateErr := s.withdrawalRepo.Update(ctx, withdrawal); updateErr != nil {
			s.logger.WithError(updateErr).Errorf("[WithdrawalService] Failed to update withdrawal status to failed")
		}

		// 4.2 退款（解冻余额）
		totalAmount := withdrawal.Amount + withdrawal.Fee
		if refundErr := s.billingAccountRepo.UpdateBalance(ctx, withdrawal.TenantID, totalAmount); refundErr != nil {
			s.logger.WithError(refundErr).Errorf("[WithdrawalService] Failed to refund balance")
		}

		// 4.3 创建退款流水
		account, _ := s.billingAccountRepo.GetByTenant(ctx, withdrawal.TenantID)
		record := &entity.WithdrawalRecord{
			TenantID:         withdrawal.TenantID,
			WithdrawalID:     withdrawal.ID,
			WithdrawalNumber: withdrawal.WithdrawalNumber,
			TransactionType:  entity.WithdrawalTransactionTypeRefund,
			Amount:           totalAmount,
			BalanceBefore:    account.AvailableCredit - totalAmount,
			BalanceAfter:     account.AvailableCredit,
			Currency:         withdrawal.Currency,
			Description:      fmt.Sprintf("提现失败，退款：提现单号=%s, 原因=%s", withdrawal.WithdrawalNumber, reason),
			TransactionDate:  now.Format("2006-01-02"),
		}
		_ = s.recordRepo.Create(ctx, record)

		return fmt.Errorf("failed to execute withdrawal: %w", err)
	}

	// 5. 更新为成功状态
	withdrawal.Status = entity.WithdrawalStatusSuccess
	withdrawal.TransactionID = &transactionID
	successTime := time.Now()
	withdrawal.SuccessAt = &successTime
	if err := s.withdrawalRepo.Update(ctx, withdrawal); err != nil {
		return fmt.Errorf("failed to update withdrawal status to success: %w", err)
	}

	// 6. 创建提现流水
	account, _ := s.billingAccountRepo.GetByTenant(ctx, withdrawal.TenantID)
	record := &entity.WithdrawalRecord{
		TenantID:         withdrawal.TenantID,
		WithdrawalID:     withdrawal.ID,
		WithdrawalNumber: withdrawal.WithdrawalNumber,
		TransactionType:  entity.WithdrawalTransactionTypeWithdraw,
		Amount:           -withdrawal.Amount,
		BalanceBefore:    account.AvailableCredit,
		BalanceAfter:     account.AvailableCredit,
		Currency:         withdrawal.Currency,
		Description:      fmt.Sprintf("提现成功：提现单号=%s, 交易ID=%s", withdrawal.WithdrawalNumber, transactionID),
		TransactionDate:  now.Format("2006-01-02"),
	}
	_ = s.recordRepo.Create(ctx, record)

	// 7. 记录日志
	s.logger.WithFields(logrus.Fields{
		"withdrawal_id":    withdrawalID,
		"transaction_id":   transactionID,
		"amount":           withdrawal.Amount,
		"execution_time_s": time.Since(startTime).Seconds(),
	}).Info("[WithdrawalService] Withdrawal executed successfully")

	// 8. 更新Prometheus指标
	UpdateWithdrawalMetrics(withdrawal.TenantID, withdrawal.WithdrawalMethod, "success", withdrawal.Amount)

	return nil
}

// GetWithdrawalStatistics 获取提现统计
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - startDate: 开始日期
//   - endDate: 结束日期
// 返回: 提现统计或错误
func (s *WithdrawalService) GetWithdrawalStatistics(
	ctx context.Context,
	tenantID string,
	startDate, endDate string,
) (*WithdrawalStatisticsResponse, error) {
	stats, err := s.withdrawalRepo.GetStatisticsByTenant(ctx, tenantID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawal statistics: %w", err)
	}

	// 计算平均金额
	var averageAmount float64
	if stats.TotalCount > 0 {
		averageAmount = stats.TotalAmount / float64(stats.TotalCount)
	}

	// 计算成功率
	var successRate float64
	if stats.TotalCount > 0 {
		successRate = float64(stats.SuccessCount) / float64(stats.TotalCount) * 100
	}

	return &WithdrawalStatisticsResponse{
		TenantID:      tenantID,
		Period: &Period{
			StartDate: startDate,
			EndDate:   endDate,
		},
		TotalCount:    stats.TotalCount,
		TotalAmount:   stats.TotalAmount,
		SuccessCount:  stats.SuccessCount,
		SuccessAmount: stats.SuccessAmount,
		PendingCount:  stats.PendingCount,
		PendingAmount: stats.PendingAmount,
		FailedCount:   stats.FailedCount,
		FailedAmount:  stats.FailedAmount,
		TotalFee:      stats.TotalFee,
		AverageAmount: averageAmount,
		SuccessRate:   successRate,
	}, nil
}

// ============================================================
// 私有辅助方法
// ============================================================

// validateWithdrawalLimit 验证提现限额
func (s *WithdrawalService) validateWithdrawalLimit(
	ctx context.Context,
	tenantID string,
	amount float64,
	limit *entity.WithdrawalLimit,
) error {
	if limit == nil {
		return nil
	}

	// 检查最小/最大金额
	if amount < limit.MinAmount {
		return fmt.Errorf("amount below minimum: %.2f < %.2f", amount, limit.MinAmount)
	}
	if amount > limit.MaxAmount {
		return fmt.Errorf("amount exceeds maximum: %.2f > %.2f", amount, limit.MaxAmount)
	}

	// 检查每日限额
	if limit.DailyLimit != nil {
		today := time.Now().Format("2006-01-02")
		todayAmount, err := s.withdrawalRepo.SumAmountByTenantAndDate(ctx, tenantID, today)
		if err != nil {
			return fmt.Errorf("failed to check daily limit: %w", err)
		}
		if todayAmount+amount > *limit.DailyLimit {
			return fmt.Errorf("daily limit exceeded: %.2f + %.2f > %.2f",
				todayAmount, amount, *limit.DailyLimit)
		}
	}

	return nil
}

// calculateFee 计算手续费
func (s *WithdrawalService) calculateFee(amount float64, limit *entity.WithdrawalLimit) float64 {
	if limit == nil {
		return 0
	}

	var fee float64
	switch limit.FeeType {
	case entity.FeeTypeFixed:
		if limit.FeeAmount != nil {
			fee = *limit.FeeAmount
		}
	case entity.FeeTypePercentage:
		if limit.FeeRate != nil {
			fee = amount * (*limit.FeeRate) / 10000
		}
	case entity.FeeTypeTiered:
		// TODO: 实现阶梯费率
		fee = 0
	}

	// 应用最低/最高手续费
	if limit.MinFee != nil && fee < *limit.MinFee {
		fee = *limit.MinFee
	}
	if limit.MaxFee != nil && fee > *limit.MaxFee {
		fee = *limit.MaxFee
	}

	// 保留两位小数
	return math.Round(fee*100) / 100
}

// performRiskCheck 执行风控检查
func (s *WithdrawalService) performRiskCheck(
	ctx context.Context,
	req *CreateWithdrawalRequest,
	account *entity.BillingAccount,
) (float64, string) {
	// 简化版风控检查
	riskScore := 0.0
	var riskReasons []string

	// 1. 检查金额是否异常
	if req.Amount > 100000 {
		riskScore += 30
		riskReasons = append(riskReasons, "大额提现")
	}

	// 2. 检查账户状态
	if account.Status != entity.BillingAccountStatusActive {
		riskScore += 50
		riskReasons = append(riskReasons, "账户状态异常")
	}

	// 3. 检查提现频率
	today := time.Now().Format("2006-01-02")
	count, err := s.withdrawalRepo.CountByTenantAndDate(ctx, req.TenantID, today)
	if err == nil && count > 5 {
		riskScore += 20
		riskReasons = append(riskReasons, "高频提现")
	}

	reason := strings.Join(riskReasons, "; ")
	if reason == "" {
		reason = "正常"
	}

	return riskScore, reason
}

// encryptBankAccount 加密银行账号
func (s *WithdrawalService) encryptBankAccount(account string) (string, error) {
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(account), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptBankAccount 解密银行账号
func (s *WithdrawalService) decryptBankAccount(encryptedAccount string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedAccount)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// generateWithdrawalNumber 生成提现单号
func (s *WithdrawalService) generateWithdrawalNumber(
	ctx context.Context,
	tenantID string,
) (string, error) {
	// 格式: WTD-{TenantID}-{YYYYMMDDHHMMSS}-{序号}
	now := time.Now()
	prefix := fmt.Sprintf("WTD-%s-%s", tenantID, now.Format("20060102150405"))

	// 查询当前秒内已有的提现数量
	count, err := s.withdrawalRepo.CountByTenantAndDate(ctx, tenantID, now.Format("2006-01-02"))
	if err != nil {
		return "", fmt.Errorf("failed to count withdrawals: %w", err)
	}

	sequence := count + 1
	withdrawalNumber := fmt.Sprintf("%s-%03d", prefix, sequence)

	return withdrawalNumber, nil
}

// getProvider 根据提现方式获取提供商
func (s *WithdrawalService) getProvider(withdrawalMethod string) string {
	switch withdrawalMethod {
	case entity.WithdrawalMethodAlipay:
		return entity.WithdrawalProviderAlipay
	case entity.WithdrawalMethodWechat:
		return entity.WithdrawalProviderWechat
	case entity.WithdrawalMethodBankTransfer:
		return entity.WithdrawalProviderBank
	default:
		return "unknown"
	}
}

// calculateEstimatedArrivalTime 计算预计到账时间
func (s *WithdrawalService) calculateEstimatedArrivalTime(withdrawal *entity.WithdrawalRequest) int64 {
	// T+1：次日到账（银行转账）
	// T+0：实时到账（支付宝、微信）
	var arrivalTime time.Time
	if withdrawal.WithdrawalMethod == entity.WithdrawalMethodBankTransfer {
		arrivalTime = time.Now().Add(24 * time.Hour)
	} else {
		arrivalTime = time.Now().Add(2 * time.Hour)
	}
	return arrivalTime.Unix()
}

// executeWithdrawalViaProvider 通过第三方支付提供商执行提现
func (s *WithdrawalService) executeWithdrawalViaProvider(
	ctx context.Context,
	withdrawal *entity.WithdrawalRequest,
) (string, error) {
	// TODO: 实现真实的第三方支付接口调用
	// 1. 根据提现方式调用不同的API
	// 2. 返回交易ID

	s.logger.Infof("[WithdrawalService] Executing withdrawal via provider: %s", withdrawal.Provider)

	// 模拟交易ID
	transactionID := fmt.Sprintf("TXN-%s-%d", withdrawal.WithdrawalNumber, time.Now().Unix())

	// 这里应该调用实际的第三方API
	// switch withdrawal.WithdrawalMethod {
	// case entity.WithdrawalMethodAlipay:
	//     return s.executeAlipayWithdrawal(ctx, withdrawal)
	// case entity.WithdrawalMethodWechat:
	//     return s.executeWechatWithdrawal(ctx, withdrawal)
	// case entity.WithdrawalMethodBankTransfer:
	//     return s.executeBankTransfer(ctx, withdrawal)
	// }

	return transactionID, nil
}

// ============================================================
// 支付宝提现
// ============================================================

// executeAlipayWithdrawal 执行支付宝提现
func (s *WithdrawalService) executeAlipayWithdrawal(
	ctx context.Context,
	withdrawal *entity.WithdrawalRequest,
) (string, error) {
	// TODO: 实现真实的支付宝提现API调用
	// 1. 构建请求参数
	// 2. 生成签名
	// 3. 发送请求
	// 4. 解析响应

	s.logger.Infof("[WithdrawalService] Executing Alipay withdrawal: %s", withdrawal.WithdrawalNumber)

	// 模拟交易ID
	transactionID := fmt.Sprintf("ALIPAY-%s-%d", withdrawal.WithdrawalNumber, time.Now().UnixNano())

	return transactionID, nil
}

// ============================================================
// 微信支付提现
// ============================================================

// executeWechatWithdrawal 执行微信支付提现
func (s *WithdrawalService) executeWechatWithdrawal(
	ctx context.Context,
	withdrawal *entity.WithdrawalRequest,
) (string, error) {
	// TODO: 实现真实的微信支付提现API调用
	s.logger.Infof("[WithdrawalService] Executing WeChat withdrawal: %s", withdrawal.WithdrawalNumber)

	// 模拟交易ID
	transactionID := fmt.Sprintf("WECHAT-%s-%d", withdrawal.WithdrawalNumber, time.Now().UnixNano())

	return transactionID, nil
}

// ============================================================
// 银行转账
// ============================================================

// executeBankTransfer 执行银行转账
func (s *WithdrawalService) executeBankTransfer(
	ctx context.Context,
	withdrawal *entity.WithdrawalRequest,
) (string, error) {
	// TODO: 实现真实的银行转账API调用
	s.logger.Infof("[WithdrawalService] Executing bank transfer: %s", withdrawal.WithdrawalNumber)

	// 模拟交易ID
	transactionID := fmt.Sprintf("BANK-%s-%d", withdrawal.WithdrawalNumber, time.Now().UnixNano())

	return transactionID, nil
}

// ============================================================
// 对账功能
// ============================================================

// Reconciliation 执行对账
func (s *WithdrawalService) Reconciliation(
	ctx context.Context,
	reconciliationDate string,
) (*entity.WithdrawalReconciliation, error) {
	s.logger.WithField("date", reconciliationDate).
		Info("[WithdrawalService] Starting reconciliation")

	startTime := time.Now()
	batch := fmt.Sprintf("RECONCILE-%s-%d", reconciliationDate, time.Now().Unix())

	// TODO: 实现真实的对账逻辑
	// 1. 查询所有提现记录
	// 2. 调用第三方对账接口
	// 3. 比对数据
	// 4. 记录差异

	reconciliation := &entity.WithdrawalReconciliation{
		ReconciliationBatch: batch,
		ReconciliationDate:  reconciliationDate,
		StartDate:           reconciliationDate,
		EndDate:             reconciliationDate,
		TotalCount:          0,
		TotalAmount:         0,
		SuccessCount:        0,
		SuccessAmount:       0,
		FailedCount:         0,
		FailedAmount:        0,
		PendingCount:        0,
		PendingAmount:       0,
		Status:              entity.ReconciliationStatusCompleted,
		DifferenceCount:     0,
		ExecutedAt:          &startTime,
		ExecutedBy:          "system",
	}
	duration := time.Since(startTime).Seconds()
	reconciliation.Duration = &duration

	if err := s.reconciliationRepo.Create(ctx, reconciliation); err != nil {
		return nil, fmt.Errorf("failed to create reconciliation: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"batch":          batch,
		"duration_s":     duration,
	}).Info("[WithdrawalService] Reconciliation completed")

	return reconciliation, nil
}

// generateBatchNo 生成批次号
func (s *WithdrawalService) generateBatchNo() string {
	return fmt.Sprintf("BATCH-%s-%d", time.Now().Format("20060102150405"), time.Now().UnixNano())
}

// generateSignature 生成签名
func (s *WithdrawalService) generateSignature(params map[string]string, secret string) string {
	// 1. 参数排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. 拼接参数
	var builder strings.Builder
	for _, k := range keys {
		if params[k] != "" {
			if builder.Len() > 0 {
				builder.WriteString("&")
			}
			builder.WriteString(k)
			builder.WriteString("=")
			builder.WriteString(params[k])
		}
	}

	// 3. 添加密钥
	builder.WriteString("&key=")
	builder.WriteString(secret)

	// 4. MD5签名
	signData := builder.String()
	hash := md5.Sum([]byte(signData))
	return fmt.Sprintf("%X", hash)
}

// parseCallback 解析回调数据
func (s *WithdrawalService) parseCallback(data map[string]interface{}) (map[string]interface{}, error) {
	// TODO: 实现真实的回调解析
	return data, nil
}

// verifyCallbackSignature 验证回调签名
func (s *WithdrawalService) verifyCallbackSignature(params map[string]interface{}, secret string) bool {
	// TODO: 实现真实的签名验证
	return true
}
