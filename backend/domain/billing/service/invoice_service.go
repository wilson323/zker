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
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gopkg.in/gomail.v2"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	billingRepository "github.com/coze-dev/coze-studio/backend/domain/billing/repository"
)

// ============================================================
// 发票管理服务 - 基于设计文档: 23-MultiTenant SaaS核心_租户计费系统.md
// ============================================================

// InvoiceService 发票服务
type InvoiceService struct {
	invoiceRepo       billingRepository.InvoiceRepository
	billingAccountRepo billingRepository.BillingAccountRepository
	db                *gorm.DB
	logger            *logrus.Logger

	// PDF生成配置
	pdfStorageDir     string
	pdfBaseURL        string

	// 邮件配置
	smtpHost          string
	smtpPort          int
	smtpUsername      string
	smtpPassword      string
	smtpFrom          string
}

// InvoiceEmailData 发票邮件数据
type InvoiceEmailData struct {
	TenantID       string
	InvoiceNumber  string
	PeriodStart    time.Time
	PeriodEnd      time.Time
	TotalAmount    float64
	Currency       string
	DueDate        time.Time
	PaymentURL     string
	CompanyName    *string
	ContactEmail   *string
}

// NewInvoiceService 创建发票服务实例
func NewInvoiceService(
	invoiceRepo billingRepository.InvoiceRepository,
	billingAccountRepo billingRepository.BillingAccountRepository,
	db *gorm.DB,
	logger *logrus.Logger,
	pdfStorageDir string,
	pdfBaseURL string,
	smtpHost string,
	smtpPort int,
	smtpUsername, smtpPassword, smtpFrom string,
) *InvoiceService {
	return &InvoiceService{
		invoiceRepo:        invoiceRepo,
		billingAccountRepo: billingAccountRepo,
		db:                 db,
		logger:             logger,
		pdfStorageDir:      pdfStorageDir,
		pdfBaseURL:         pdfBaseURL,
		smtpHost:           smtpHost,
		smtpPort:           smtpPort,
		smtpUsername:       smtpUsername,
		smtpPassword:       smtpPassword,
		smtpFrom:           smtpFrom,
	}
}

// ============================================================
// 核心方法实现
// ============================================================

// GenerateInvoicePDF 生成发票PDF
// 参数:
//   - ctx: 上下文
//   - invoiceID: 发票ID
// 返回: PDF文件路径或错误
func (s *InvoiceService) GenerateInvoicePDF(
	ctx context.Context,
	invoiceID uint64,
) (string, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		RecordInvoiceOperation("generate_pdf", duration)
	}()

	s.logger.WithFields(logrus.Fields{
		"invoice_id": invoiceID,
	}).Info("[InvoiceService] Generating invoice PDF")

	// 1. 查询发票
	invoice, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		s.logger.WithError(err).Errorf("[InvoiceService] Failed to get invoice")
		return "", fmt.Errorf("failed to get invoice: %w", err)
	}

	if invoice == nil {
		return "", fmt.Errorf("invoice not found: %d", invoiceID)
	}

	// 2. 查询发票明细
	lineItems, err := s.invoiceRepo.GetLineItems(ctx, invoiceID)
	if err != nil {
		s.logger.WithError(err).Errorf("[InvoiceService] Failed to get line items")
		return "", fmt.Errorf("failed to get line items: %w", err)
	}

	// 3. 查询计费账户信息
	account, err := s.billingAccountRepo.GetByID(ctx, invoice.BillingAccountID)
	if err != nil {
		s.logger.WithError(err).Errorf("[InvoiceService] Failed to get billing account")
		return "", fmt.Errorf("failed to get billing account: %w", err)
	}

	// 4. 生成PDF文件
	pdfPath := s.getPDFPath(invoice.InvoiceNumber)

	if err := s.createInvoicePDF(invoice, lineItems, account, pdfPath); err != nil {
		s.logger.WithError(err).Errorf("[InvoiceService] Failed to create PDF")
		return "", fmt.Errorf("failed to create PDF: %w", err)
	}

	// 5. 更新发票记录
	pdfURL := s.pdfBaseURL + "/" + filepath.Base(pdfPath)
	now := time.Now()
	invoice.PdfURL = &pdfURL
	invoice.PdfGeneratedAt = &now

	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		s.logger.WithError(err).Errorf("[InvoiceService] Failed to update invoice")
		return "", fmt.Errorf("failed to update invoice: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"invoice_id":        invoiceID,
		"invoice_number":    invoice.InvoiceNumber,
		"pdf_path":          pdfPath,
		"pdf_url":           pdfURL,
		"generation_time_s": time.Since(startTime).Seconds(),
	}).Info("[InvoiceService] Invoice PDF generated successfully")

	return pdfPath, nil
}

// SendInvoiceEmail 发送发票邮件
// 参数:
//   - ctx: 上下文
//   - invoiceID: 发票ID
//   - recipientEmail: 收件人邮箱
// 返回: 错误
func (s *InvoiceService) SendInvoiceEmail(
	ctx context.Context,
	invoiceID uint64,
	recipientEmail string,
) error {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		RecordInvoiceOperation("send_email", duration)
	}()

	s.logger.WithFields(logrus.Fields{
		"invoice_id":     invoiceID,
		"recipient_email": recipientEmail,
	}).Info("[InvoiceService] Sending invoice email")

	// 1. 查询发票
	invoice, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		s.logger.WithError(err).Errorf("[InvoiceService] Failed to get invoice")
		return fmt.Errorf("failed to get invoice: %w", err)
	}

	if invoice == nil {
		return fmt.Errorf("invoice not found: %d", invoiceID)
	}

	// 2. 检查PDF是否已生成
	if invoice.PdfURL == nil {
		// 自动生成PDF
		if _, err := s.GenerateInvoicePDF(ctx, invoiceID); err != nil {
			s.logger.WithError(err).Warnf("[InvoiceService] Failed to generate PDF for email")
			return fmt.Errorf("failed to generate PDF: %w", err)
		}

		// 重新查询发票以获取PDF URL
		invoice, _ = s.invoiceRepo.GetByID(ctx, invoiceID)
	}

	// 3. 构建邮件数据
	emailData := &InvoiceEmailData{
		TenantID:      invoice.TenantID,
		InvoiceNumber: invoice.InvoiceNumber,
		PeriodStart:   invoice.PeriodStart,
		PeriodEnd:     invoice.PeriodEnd,
		TotalAmount:   invoice.TotalAmount,
		Currency:      invoice.Currency,
		PaymentURL:    fmt.Sprintf("%s/payments/%d", s.pdfBaseURL, invoiceID),
	}

	if invoice.DueDate != nil {
		emailData.DueDate = *invoice.DueDate
	}

	// 4. 构建邮件内容
	subject := fmt.Sprintf("Invoice %s - %.2f %s", invoice.InvoiceNumber, invoice.TotalAmount, invoice.Currency)
	body := s.buildInvoiceEmailBody(emailData)

	// 5. 创建邮件消息
	mail := gomail.NewMessage()
	mail.SetHeader("From", s.smtpFrom)
	mail.SetHeader("To", recipientEmail)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", body)

	// 附加PDF文件
	if invoice.PdfURL != nil {
		pdfPath := s.getPDFPath(invoice.InvoiceNumber)
		if _, err := os.Stat(pdfPath); err == nil {
			mail.Attach(pdfPath)
		}
	}

	// 6. 发送邮件
	dialer := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUsername, s.smtpPassword)

	if err := dialer.DialAndSend(mail); err != nil {
		s.logger.WithError(err).Errorf("[InvoiceService] Failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	// 7. 更新发票状态（如果还未发送）
	if invoice.Status == entity.InvoiceStatusDraft {
		invoice.Status = entity.InvoiceStatusSent
		if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
			s.logger.WithError(err).Warnf("[InvoiceService] Failed to update invoice status")
		}
	}

	s.logger.WithFields(logrus.Fields{
		"invoice_id":    invoiceID,
		"invoice_number": invoice.InvoiceNumber,
		"recipient":     recipientEmail,
		"send_time_s":   time.Since(startTime).Seconds(),
	}).Info("[InvoiceService] Invoice email sent successfully")

	return nil
}

// GetInvoiceWithDetails 获取发票及明细
// 参数:
//   - ctx: 上下文
//   - invoiceID: 发票ID
// 返回: 发票及明细或错误
func (s *InvoiceService) GetInvoiceWithDetails(
	ctx context.Context,
	invoiceID uint64,
) (*entity.Invoice, []entity.InvoiceLineItem, error) {
	invoice, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	if invoice == nil {
		return nil, nil, fmt.Errorf("invoice not found: %d", invoiceID)
	}

	lineItems, err := s.invoiceRepo.GetLineItems(ctx, invoiceID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get line items: %w", err)
	}

	return invoice, lineItems, nil
}

// ListInvoicesByTenant 列出租户的发票
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - page: 页码
//   - pageSize: 每页大小
// 返回: 发票列表、总数或错误
func (s *InvoiceService) ListInvoicesByTenant(
	ctx context.Context,
	tenantID string,
	page, pageSize int,
) ([]*entity.Invoice, int64, error) {
	invoices, total, err := s.invoiceRepo.ListByTenant(ctx, tenantID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list invoices: %w", err)
	}

	return invoices, total, nil
}

// ============================================================
// 私有辅助方法
// ============================================================

// getPDFPath 获取PDF文件路径
func (s *InvoiceService) getPDFPath(invoiceNumber string) string {
	// 确保存储目录存在
	os.MkdirAll(s.pdfStorageDir, 0755)

	filename := fmt.Sprintf("%s.pdf", invoiceNumber)
	return filepath.Join(s.pdfStorageDir, filename)
}

// createInvoicePDF 创建发票PDF文件
func (s *InvoiceService) createInvoicePDF(
	invoice *entity.Invoice,
	lineItems []entity.InvoiceLineItem,
	account *entity.BillingAccount,
	pdfPath string,
) error {
	// 1. 创建PDF文档
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// 2. 设置字体
	pdf.SetFont("Arial", "", 12)

	// 3. 添加标题
	pdf.SetFontSize(20)
	pdf.Cell(190, 15, "INVOICE")
	pdf.Ln(20)

	// 4. 添加发票信息
	pdf.SetFontSize(12)
	pdf.Cell(95, 10, fmt.Sprintf("Invoice Number: %s", invoice.InvoiceNumber))
	pdf.Cell(95, 10, fmt.Sprintf("Date: %s", invoice.CreatedAt.Format("2006-01-02")))
	pdf.Ln(8)

	pdf.Cell(95, 10, fmt.Sprintf("Tenant ID: %s", invoice.TenantID))
	if invoice.DueDate != nil {
		pdf.Cell(95, 10, fmt.Sprintf("Due Date: %s", invoice.DueDate.Format("2006-01-02")))
	}
	pdf.Ln(10)

	pdf.Cell(95, 10, fmt.Sprintf("Period: %s to %s",
		invoice.PeriodStart.Format("2006-01-02"),
		invoice.PeriodEnd.Format("2006-01-02")))
	pdf.Ln(15)

	// 5. 添加费用明细表头
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(20, 10, "#", "1", 0, "C", true, 0, "")
	pdf.CellFormat(90, 10, "Description", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 10, "Qty", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Unit Price", "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 10, "Amount", "1", 1, "R", true, 0, "")
	pdf.Ln(5)

	// 6. 添加费用明细
	pdf.SetFillColor(240, 240, 240)
	for i, item := range lineItems {
		if i%2 == 0 {
			pdf.SetFillColor(240, 240, 240)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		pdf.CellFormat(20, 8, fmt.Sprintf("%d", i+1), "1", 0, "C", true, 0, "")
		pdf.CellFormat(90, 8, truncateString(item.Description, 40), "1", 0, "L", true, 0, "")
		pdf.CellFormat(20, 8, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", true, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%.2f", item.UnitPrice), "1", 0, "R", true, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%.2f", item.Amount), "1", 1, "R", true, 0, "")
	}
	pdf.Ln(5)

	// 7. 添加汇总（企业级发票：符合财务规范的PDF格式）
	pdf.Cell(140, 10, "")
	pdf.CellFormat(25, 10, "Subtotal:", "", 0, "R", false, 0, "")
	pdf.CellFormat(25, 10, fmt.Sprintf("%.2f %s", invoice.Subtotal, invoice.Currency), "", 1, "R", false, 0, "")

	if invoice.TaxAmount > 0 {
		pdf.Cell(140, 8, "")
		pdf.CellFormat(25, 8, "Tax:", "", 0, "R", false, 0, "")
		pdf.CellFormat(25, 8, fmt.Sprintf("%.2f %s", invoice.TaxAmount, invoice.Currency), "", 1, "R", false, 0, "")
	}

	if invoice.DiscountAmount > 0 {
		pdf.Cell(140, 8, "")
		pdf.CellFormat(25, 8, "Discount:", "", 0, "R", false, 0, "")
		pdf.CellFormat(25, 8, fmt.Sprintf("-%.2f %s", invoice.DiscountAmount, invoice.Currency), "", 1, "R", false, 0, "")
	}

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(140, 10, "")
	pdf.CellFormat(25, 10, "Total:", "", 0, "R", false, 0, "")
	pdf.CellFormat(25, 10, fmt.Sprintf("%.2f %s", invoice.TotalAmount, invoice.Currency), "", 1, "R", false, 0, "")
	pdf.Ln(15)

	// 8. 添加支付信息
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(190, 8, "Payment Information:")
	pdf.Ln(6)
	pdf.Cell(190, 6, "Please pay by the due date to avoid service interruption.")
	pdf.Ln(6)
	pdf.Cell(190, 6, fmt.Sprintf("Payment Methods: Alipay, WeChat Pay, Credit Card"))
	pdf.Ln(10)

	// 9. 添加备注
	if invoice.Notes != nil {
		pdf.Cell(190, 8, "Notes:")
		pdf.Ln(6)
		pdf.MultiCell(190, 6, *invoice.Notes, "", "", false)
	}

	// 10. 保存PDF文件
	return pdf.OutputFileAndClose(pdfPath)
}

// buildInvoiceEmailBody 构建发票邮件正文
func (s *InvoiceService) buildInvoiceEmailBody(data *InvoiceEmailData) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .invoice-info { background-color: white; padding: 15px; margin: 15px 0; border-radius: 5px; }
        .total { font-size: 24px; font-weight: bold; color: #4CAF50; }
        .button { display: inline-block; padding: 12px 24px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px; margin: 15px 0; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #777; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Invoice Ready</h1>
        </div>
        <div class="content">
            <p>Hello,</p>
            <p>Your invoice is now ready for the billing period %s to %s.</p>

            <div class="invoice-info">
                <h2>Invoice Details</h2>
                <p><strong>Invoice Number:</strong> %s</p>
                <p><strong>Billing Period:</strong> %s to %s</p>
                <p><strong>Due Date:</strong> %s</p>
                <p class="total">Total Amount: %.2f %s</p>
            </div>

            <p>Please complete your payment by the due date to avoid service interruption.</p>
            <p style="text-align: center;">
                <a href="%s" class="button">Pay Now</a>
            </p>
            <p>Or visit your account dashboard to view payment options.</p>
        </div>
        <div class="footer">
            <p>Thank you for your business!</p>
            <p>If you have any questions, please contact our support team.</p>
        </div>
    </div>
</body>
</html>
`,
		data.PeriodStart.Format("2006-01-02"),
		data.PeriodEnd.Format("2006-01-02"),
		data.InvoiceNumber,
		data.PeriodStart.Format("2006-01-02"),
		data.PeriodEnd.Format("2006-01-02"),
		data.DueDate.Format("2006-01-02"),
		data.TotalAmount,
		data.Currency,
		data.PaymentURL,
	)
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
