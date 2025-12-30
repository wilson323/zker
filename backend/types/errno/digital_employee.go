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

package errno

import (
	"net/http"
)

// 数字员工错误码 (3xxx系列)
var (
	// 员工画像错误 (30xxx)
	EmployeeNotFound = &BaseErrorCode{
		code:       "DIGITAL_EMPLOYEE_30001",
		message:    "Employee not found",
		messageZH:  "员工不存在",
		messageEN:  "Employee not found",
		httpStatus: http.StatusNotFound,
	}
	EmployeeNameExists = &BaseErrorCode{
		code:       "DIGITAL_EMPLOYEE_30002",
		message:    "Employee name already exists",
		messageZH:  "员工名称已存在",
		messageEN:  "Employee name already exists",
		httpStatus: http.StatusConflict,
	}
	EmployeeInvalidRole = &BaseErrorCode{
		code:       "DIGITAL_EMPLOYEE_30003",
		message:    "Invalid employee role",
		messageZH:  "无效的员工角色",
		messageEN:  "Invalid employee role",
		httpStatus: http.StatusBadRequest,
	}
	EmployeeInactive = &BaseErrorCode{
		code:       "DIGITAL_EMPLOYEE_30004",
		message:    "Employee is inactive",
		messageZH:  "员工未激活",
		messageEN:  "Employee is inactive",
		httpStatus: http.StatusBadRequest,
	}
	EmployeeSkillMissing = &BaseErrorCode{
		code:       "DIGITAL_EMPLOYEE_30005",
		message:    "Required skills missing",
		messageZH:  "缺少必需技能",
		messageEN:  "Required skills missing",
		httpStatus: http.StatusBadRequest,
	}

	// 任务分配错误 (31xxx)
	TaskNotFound = &BaseErrorCode{
		code:       "DIGITAL_EMPLOYEE_31001",
		message:    "Task assignment not found",
		messageZH:  "任务分配不存在",
		messageEN:  "Task assignment not found",
		httpStatus: http.StatusNotFound,
	}
	TaskAlreadyAssigned = &BaseErrorCode{
		code:       "DIGITAL_EMPLOYEE_31002",
		message:    "Task already assigned",
		messageZH:  "任务已分配",
		messageEN:  "Task already assigned",
		httpStatus: http.StatusConflict,
	}
	TaskInvalidStatus = &BaseErrorCode{
		code:       "DIGITAL_EMPLOYEE_31003",
		message:    "Invalid task status",
		messageZH:  "无效的任务状态",
		messageEN:  "Invalid task status",
		httpStatus: http.StatusBadRequest,
	}
	NoAvailableEmployee = &BaseErrorCode{
		code:       "DIGITAL_EMPLOYEE_31004",
		message:    "No available employee found",
		messageZH:  "未找到可用员工",
		messageEN:  "No available employee found",
		httpStatus: http.StatusNotFound,
	}
	TaskAlreadyCompleted = &BaseErrorCode{
		code:       "DIGITAL_EMPLOYEE_31005",
		message:    "Task already completed",
		messageZH:  "任务已完成",
		messageEN:  "Task already completed",
		httpStatus: http.StatusBadRequest,
	}

	// 绩效统计错误 (32xxx)
	PerformanceNotFound = &BaseErrorCode{
		code:       "DIGITAL_EMPLOYEE_32001",
		message:    "Performance record not found",
		messageZH:  "绩效记录不存在",
		messageEN:  "Performance record not found",
		httpStatus: http.StatusNotFound,
	}
	PerformanceInvalidPeriod = &BaseErrorCode{
		code:       "DIGITAL_EMPLOYEE_32002",
		message:    "Invalid performance period",
		messageZH:  "无效的统计周期",
		messageEN:  "Invalid performance period",
		httpStatus: http.StatusBadRequest,
	}
	PerformanceCalculationFailed = &BaseErrorCode{
		code:       "DIGITAL_EMPLOYEE_32003",
		message:    "Failed to calculate performance",
		messageZH:  "绩效计算失败",
		messageEN:  "Failed to calculate performance",
		httpStatus: http.StatusInternalServerError,
	}
)
