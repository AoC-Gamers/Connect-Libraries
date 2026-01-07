/**
 * TypeScript Types for Connect Backend Error Responses
 * 
 * Copy this file to your frontend project:
 * src/models/api-error.model.ts
 * 
 * These types match the Go ErrorResponse struct from connect-errors library
 */

/**
 * Standardized error codes from backend
 * Must match codes.go ErrorCode constants
 */
export enum ErrorCode {
  // Authentication & Authorization
  UNAUTHORIZED = 'UNAUTHORIZED',
  TOKEN_EXPIRED = 'TOKEN_EXPIRED',
  TOKEN_INVALID = 'TOKEN_INVALID',
  POLICY_VERSION_MISMATCH = 'POLICY_VERSION_MISMATCH',
  PERMISSION_DENIED = 'PERMISSION_DENIED',
  INSUFFICIENT_PERMISSIONS = 'INSUFFICIENT_PERMISSIONS',
  
  // Resource Errors
  NOT_FOUND = 'NOT_FOUND',
  MEMBERSHIP_NOT_FOUND = 'MEMBERSHIP_NOT_FOUND',
  ALREADY_EXISTS = 'ALREADY_EXISTS',
  CONFLICT = 'CONFLICT',
  RESOURCE_LOCKED = 'RESOURCE_LOCKED',
  
  // Validation Errors
  VALIDATION_ERROR = 'VALIDATION_ERROR',
  INVALID_REQUEST = 'INVALID_REQUEST',
  BAD_REQUEST = 'BAD_REQUEST',
  MISSING_REQUIRED_FIELD = 'MISSING_REQUIRED_FIELD',
  INVALID_FORMAT = 'INVALID_FORMAT',
  OUT_OF_RANGE = 'OUT_OF_RANGE',
  
  // Server Errors
  INTERNAL_ERROR = 'INTERNAL_ERROR',
  DATABASE_ERROR = 'DATABASE_ERROR',
  SERVICE_UNAVAILABLE = 'SERVICE_UNAVAILABLE',
  TIMEOUT = 'TIMEOUT',
  
  // Business Logic Errors
  OPERATION_NOT_ALLOWED = 'OPERATION_NOT_ALLOWED',
  QUOTA_EXCEEDED = 'QUOTA_EXCEEDED',
  RATE_LIMIT_EXCEEDED = 'RATE_LIMIT_EXCEEDED',
}

/**
 * Error response structure from backend
 * Matches ErrorResponse struct from errors.go
 */
export interface ApiErrorResponse {
  /** Short error message (legacy compatibility) */
  error: string;
  
  /** Programmatic error code */
  code?: ErrorCode;
  
  /** HTTP status code */
  status: number;
  
  /** Detailed explanation of the error */
  detail?: string;
  
  /** Additional metadata specific to the error context */
  meta?: Record<string, unknown>;
}

/**
 * Extended Error class with structured information
 */
export class ApiError extends Error {
  public readonly code?: ErrorCode;
  public readonly status: number;
  public readonly detail?: string;
  public readonly meta?: Record<string, unknown>;

  constructor(response: ApiErrorResponse) {
    super(response.error);
    this.name = 'ApiError';
    this.code = response.code;
    this.status = response.status;
    this.detail = response.detail;
    this.meta = response.meta;
  }

  /**
   * Check if error is of specific code
   */
  is(code: ErrorCode): boolean {
    return this.code === code;
  }

  /**
   * Check if error is permission-related
   */
  isPermissionError(): boolean {
    return this.code === ErrorCode.PERMISSION_DENIED || 
           this.code === ErrorCode.INSUFFICIENT_PERMISSIONS;
  }

  /**
   * Check if error is membership not found
   */
  isMembershipNotFound(): boolean {
    return this.code === ErrorCode.MEMBERSHIP_NOT_FOUND;
  }

  /**
   * Check if error is not found
   */
  isNotFound(): boolean {
    return this.code === ErrorCode.NOT_FOUND || 
           this.code === ErrorCode.MEMBERSHIP_NOT_FOUND;
  }

  /**
   * Check if error requires re-authentication
   */
  requiresReauth(): boolean {
    return this.code === ErrorCode.TOKEN_EXPIRED || 
           this.code === ErrorCode.POLICY_VERSION_MISMATCH ||
           this.code === ErrorCode.TOKEN_INVALID;
  }

  /**
   * Check if error is validation-related
   */
  isValidationError(): boolean {
    return this.code === ErrorCode.VALIDATION_ERROR ||
           this.code === ErrorCode.INVALID_FORMAT ||
           this.code === ErrorCode.MISSING_REQUIRED_FIELD ||
           this.code === ErrorCode.OUT_OF_RANGE;
  }

  /**
   * Check if error is server-side
   */
  isServerError(): boolean {
    return this.status >= 500 && this.status < 600;
  }

  /**
   * Get user-friendly message
   */
  getUserMessage(): string {
    const friendlyMessages: Record<ErrorCode, string> = {
      [ErrorCode.UNAUTHORIZED]: 'You need to be logged in',
      [ErrorCode.TOKEN_EXPIRED]: 'Your session has expired',
      [ErrorCode.TOKEN_INVALID]: 'Your session is invalid',
      [ErrorCode.POLICY_VERSION_MISMATCH]: 'Please log in again',
      [ErrorCode.PERMISSION_DENIED]: 'You don\'t have permission for this action',
      [ErrorCode.INSUFFICIENT_PERMISSIONS]: 'Insufficient permissions',
      [ErrorCode.NOT_FOUND]: 'Resource not found',
      [ErrorCode.MEMBERSHIP_NOT_FOUND]: 'You don\'t have access to this resource',
      [ErrorCode.ALREADY_EXISTS]: 'Resource already exists',
      [ErrorCode.CONFLICT]: 'Conflict with current state',
      [ErrorCode.RESOURCE_LOCKED]: 'Resource is locked',
      [ErrorCode.VALIDATION_ERROR]: 'Invalid data provided',
      [ErrorCode.INVALID_REQUEST]: 'Invalid request',
      [ErrorCode.BAD_REQUEST]: 'Bad request',
      [ErrorCode.MISSING_REQUIRED_FIELD]: 'Required field is missing',
      [ErrorCode.INVALID_FORMAT]: 'Invalid format',
      [ErrorCode.OUT_OF_RANGE]: 'Value out of range',
      [ErrorCode.INTERNAL_ERROR]: 'An internal error occurred',
      [ErrorCode.DATABASE_ERROR]: 'Database error occurred',
      [ErrorCode.SERVICE_UNAVAILABLE]: 'Service temporarily unavailable',
      [ErrorCode.TIMEOUT]: 'Operation timed out',
      [ErrorCode.OPERATION_NOT_ALLOWED]: 'Operation not allowed',
      [ErrorCode.QUOTA_EXCEEDED]: 'Quota exceeded',
      [ErrorCode.RATE_LIMIT_EXCEEDED]: 'Too many requests',
    };

    return this.code && friendlyMessages[this.code] 
      ? friendlyMessages[this.code]
      : this.message || 'An error occurred';
  }

  /**
   * Convert to plain object for logging
   */
  toJSON() {
    return {
      name: this.name,
      message: this.message,
      code: this.code,
      status: this.status,
      detail: this.detail,
      meta: this.meta,
    };
  }
}

/**
 * Type guard to check if error is ApiError
 */
export function isApiError(error: unknown): error is ApiError {
  return error instanceof ApiError;
}

/**
 * Type guard to check if response has error structure
 */
export function isApiErrorResponse(data: unknown): data is ApiErrorResponse {
  return (
    typeof data === 'object' &&
    data !== null &&
    'error' in data &&
    'status' in data
  );
}
