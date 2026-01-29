// API Response wrapper
export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: ErrorInfo
}

export interface ErrorInfo {
  code: string
  message: string
}

// Environment types
export interface Environment {
  path: string
  name: string
  cve?: string[]
  app: string
  tags?: string[]
  downloaded: boolean
  running: boolean
}

export interface EnvironmentList {
  environments: Environment[]
  total: number
}

export interface EnvironmentStatus {
  environment: Environment
  containers?: ContainerStatus[]
  running: boolean
  local_path?: string
}

export interface StatusList {
  environments: EnvironmentStatus[]
  total: number
}

export interface ContainerStatus {
  id: string
  name: string
  image: string
  status: string
  state: string
  ports?: PortMapping[]
  created_at?: string
  started_at?: string
}

export interface PortMapping {
  host_ip?: string
  host_port: string
  container_port: string
  protocol: string
}

export interface EnvironmentInfo {
  environment: Environment
  readme?: string
  compose_file?: string
  downloaded: boolean
  local_path?: string
}

// Request types
export interface StartRequest {
  pull?: boolean
  build?: boolean
  force_recreate?: boolean
}

export interface CleanRequest {
  remove_volumes?: boolean
  remove_images?: boolean
  remove_files?: boolean
}

// System types
export interface SystemStatus {
  initialized: boolean
  last_sync_time?: string
  need_sync: boolean
  version: string
}

export interface SyncupResult {
  success: boolean
  last_sync_time: string
  total: number
}

// UI types
export type FilterType = 'all' | 'downloaded' | 'running'
