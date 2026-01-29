import { fetchApi, postApi, deleteApi } from './client'
import type {
  EnvironmentList,
  StatusList,
  EnvironmentInfo,
  EnvironmentStatus,
  SystemStatus,
  SyncupResult,
  StartRequest,
  CleanRequest,
} from '@/types'

// System
export async function getSystemStatus(): Promise<SystemStatus> {
  return fetchApi<SystemStatus>('status')
}

export async function syncEnvironments(): Promise<SyncupResult> {
  return postApi<SyncupResult>('syncup')
}

// Environment lists
export async function listEnvironments(): Promise<EnvironmentList> {
  return fetchApi<EnvironmentList>('environments')
}

export async function listDownloadedEnvironments(): Promise<StatusList> {
  return fetchApi<StatusList>('environments/downloaded')
}

export async function listRunningEnvironments(): Promise<StatusList> {
  return fetchApi<StatusList>('environments/running')
}

// Single environment
export async function getEnvironmentInfo(path: string): Promise<EnvironmentInfo> {
  return fetchApi<EnvironmentInfo>(`environments/info/${path}`)
}

export async function getEnvironmentStatus(path: string): Promise<EnvironmentStatus> {
  return fetchApi<EnvironmentStatus>(`environments/status/${path}`)
}

// Environment actions
export async function startEnvironment(path: string, options?: StartRequest): Promise<EnvironmentStatus> {
  return postApi<EnvironmentStatus>(`environments/start/${path}`, options)
}

export async function stopEnvironment(path: string): Promise<EnvironmentStatus> {
  return postApi<EnvironmentStatus>(`environments/stop/${path}`)
}

export async function restartEnvironment(path: string): Promise<EnvironmentStatus> {
  return postApi<EnvironmentStatus>(`environments/restart/${path}`)
}

export async function cleanEnvironment(path: string, options?: CleanRequest): Promise<void> {
  return deleteApi<void>(`environments/clean/${path}`, options)
}
