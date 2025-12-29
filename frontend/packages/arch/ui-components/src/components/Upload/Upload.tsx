// frontend/packages/arch/ui-components/src/components/Upload/Upload.tsx

import React, { useState, useRef, CSSProperties } from 'react';
import { useStyles } from './Upload.styles';

export interface UploadFile {
  uid: string;
  name: string;
  status: 'uploading' | 'done' | 'error';
  percent: number;
  url?: string;
  response?: any;
  error?: any;
}

export interface UploadProps {
  /**
   * 接受上传的文件类型
   */
  accept?: string;

  /**
   * 是否支持多文件上传
   */
  multiple?: boolean;

  /**
   * 最大上传数量
   */
  maxCount?: number;

  /**
   * 上传地址
   */
  action?: string;

  /**
   * 上传请求的 HTTP method
   */
  method?: 'POST' | 'PUT' | 'PATCH';

  /**
   * 上传请求的 headers
   */
  headers?: Record<string, string>;

  /**
   * 上传请求的额外数据
   */
  data?: Record<string, any> | ((file: UploadFile) => Record<string, any>);

  /**
   * 是否禁用
   */
  disabled?: boolean;

  /**
   * 是否支持拖拽上传
   */
  drag?: boolean;

  /**
   * 文件列表
   */
  fileList?: UploadFile[];

  /**
   * 文件变更回调
   */
  onChange?: (fileList: UploadFile[]) => void;

  /**
   * 自定义上传实现
   */
  customRequest?: (options: {
    file: UploadFile;
    onProgress: (percent: number) => void;
    onSuccess: (response: any) => void;
    onError: (error: any) => void;
  }) => void;

  /**
   * 自定义样式类名
   */
  className?: string;

  /**
   * 自定义样式
   */
  style?: CSSProperties;

  /**
   * 子元素（自定义触发按钮）
   */
  children?: React.ReactNode;
}

export const Upload: React.FC<UploadProps> = ({
  accept,
  multiple = false,
  maxCount,
  action = '/api/upload',
  method = 'POST',
  headers = {},
  data,
  disabled = false,
  drag = false,
  fileList: externalFileList,
  onChange,
  customRequest,
  className,
  style,
  children,
}) => {
  const classes = useStyles();
  const [fileList, setFileList] = useState<UploadFile[]>(externalFileList || []);
  const [isDragging, setIsDragging] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const controlled = externalFileList !== undefined;
  const currentFileList = controlled ? externalFileList : fileList;

  const triggerUpload = () => {
    inputRef.current?.click();
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(event.target.files || []);
    processFiles(files);
    // 重置 input 以允许重复上传相同文件
    event.target.value = '';
  };

  const processFiles = (files: File[]) => {
    if (maxCount && currentFileList.length + files.length > maxCount) {
      alert(`最多只能上传 ${maxCount} 个文件`);
      return;
    }

    const newFiles: UploadFile[] = files.map((file) => ({
      uid: `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      name: file.name,
      status: 'uploading',
      percent: 0,
    }));

    const updatedList = [...currentFileList, ...newFiles];

    if (!controlled) {
      setFileList(updatedList);
    }
    onChange?.(updatedList);

    // 上传文件
    newFiles.forEach((uploadFile, index) => {
      const file = files[index];
      uploadFile(uploadFile, file);
    });
  };

  const uploadFile = async (uploadFile: UploadFile, file: File) => {
    if (customRequest) {
      customRequest({
        file: uploadFile,
        onProgress: (percent) => {
          updateFileProgress(uploadFile.uid, percent);
        },
        onSuccess: (response) => {
          updateFileStatus(uploadFile.uid, 'done', response);
        },
        onError: (error) => {
          updateFileStatus(uploadFile.uid, 'error', error);
        },
      });
      return;
    }

    // 默认上传实现
    const formData = new FormData();
    formData.append('file', file);

    const requestData = typeof data === 'function' ? data(uploadFile) : data;
    if (requestData) {
      Object.keys(requestData).forEach((key) => {
        formData.append(key, requestData[key]);
      });
    }

    try {
      const xhr = new XMLHttpRequest();

      xhr.upload.onprogress = (event) => {
        if (event.lengthComputable) {
          const percent = Math.round((event.loaded / event.total) * 100);
          updateFileProgress(uploadFile.uid, percent);
        }
      };

      xhr.onload = () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          const response = JSON.parse(xhr.responseText);
          updateFileStatus(uploadFile.uid, 'done', response);
        } else {
          updateFileStatus(uploadFile.uid, 'error', xhr.responseText);
        }
      };

      xhr.onerror = () => {
        updateFileStatus(uploadFile.uid, 'error', '上传失败');
      };

      xhr.open(method, action, true);
      Object.keys(headers).forEach((key) => {
        xhr.setRequestHeader(key, headers[key]);
      });
      xhr.send(formData);
    } catch (error) {
      updateFileStatus(uploadFile.uid, 'error', error);
    }
  };

  const updateFileProgress = (uid: string, percent: number) => {
    const updatedList = currentFileList.map((file) =>
      file.uid === uid ? { ...file, percent } : file
    );
    if (!controlled) {
      setFileList(updatedList);
    }
    onChange?.(updatedList);
  };

  const updateFileStatus = (
    uid: string,
    status: 'done' | 'error',
    response?: any
  ) => {
    const updatedList = currentFileList.map((file) =>
      file.uid === uid ? { ...file, status, response, percent: status === 'done' ? 100 : file.percent } : file
    );
    if (!controlled) {
      setFileList(updatedList);
    }
    onChange?.(updatedList);
  };

  const handleRemove = (uid: string) => {
    const updatedList = currentFileList.filter((file) => file.uid !== uid);
    if (!controlled) {
      setFileList(updatedList);
    }
    onChange?.(updatedList);
  };

  const handleDragOver = (event: React.DragEvent) => {
    event.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const handleDrop = (event: React.DragEvent) => {
    event.preventDefault();
    setIsDragging(false);

    if (disabled) return;

    const files = Array.from(event.dataTransfer.files);
    processFiles(files);
  };

  const renderUploadButton = () => {
    if (children) {
      return <div onClick={!disabled ? triggerUpload : undefined}>{children}</div>;
    }

    return (
      <button
        className={classes.uploadButton}
        onClick={!disabled ? triggerUpload : undefined}
        disabled={disabled}
      >
        选择文件
      </button>
    );
  };

  return (
    <div className={`${classes.container} ${className || ''}`} style={style}>
      <input
        ref={inputRef}
        type="file"
        accept={accept}
        multiple={multiple}
        onChange={handleFileChange}
        style={{ display: 'none' }}
      />

      {drag ? (
        <div
          className={`${classes.dragger} ${isDragging ? classes.dragging : ''}`}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
        >
          <div className={classes.dragContent}>
            <div className={classes.dragIcon}>📁</div>
            <div className={classes.dragText}>
              点击或拖拽文件到此区域上传
            </div>
            {renderUploadButton()}
          </div>
        </div>
      ) : (
        renderUploadButton()
      )}

      {currentFileList.length > 0 && (
        <div className={classes.list}>
          {currentFileList.map((file) => (
            <div key={file.uid} className={classes.listItem}>
              <div className={classes.fileInfo}>
                <span className={classes.fileName}>{file.name}</span>
                <span className={classes.fileStatus}>
                  {file.status === 'uploading' && `${file.percent}%`}
                  {file.status === 'done' && '✓'}
                  {file.status === 'error' && '✗'}
                </span>
              </div>
              {file.status === 'uploading' && (
                <div className={classes.progressBar}>
                  <div
                    className={classes.progressFill}
                    style={{ width: `${file.percent}%` }}
                  />
                </div>
              )}
              <button
                className={classes.removeButton}
                onClick={() => handleRemove(file.uid)}
              >
                ✕
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Upload;
