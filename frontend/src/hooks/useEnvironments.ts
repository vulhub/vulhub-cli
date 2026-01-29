import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listEnvironments,
  listDownloadedEnvironments,
  listRunningEnvironments,
  getEnvironmentInfo,
  getEnvironmentStatus,
  getSystemStatus,
  syncEnvironments,
  startEnvironment,
  stopEnvironment,
  restartEnvironment,
  cleanEnvironment,
} from '@/api/environments'
import type { StartRequest, CleanRequest } from '@/types'
import { toast } from '@/hooks/use-toast'

// Query keys
export const queryKeys = {
  systemStatus: ['systemStatus'] as const,
  environments: ['environments'] as const,
  downloaded: ['environments', 'downloaded'] as const,
  running: ['environments', 'running'] as const,
  info: (path: string) => ['environments', 'info', path] as const,
  status: (path: string) => ['environments', 'status', path] as const,
}

// System status
export function useSystemStatus() {
  return useQuery({
    queryKey: queryKeys.systemStatus,
    queryFn: getSystemStatus,
  })
}

// Environment lists
export function useEnvironments() {
  return useQuery({
    queryKey: queryKeys.environments,
    queryFn: listEnvironments,
  })
}

export function useDownloadedEnvironments() {
  return useQuery({
    queryKey: queryKeys.downloaded,
    queryFn: listDownloadedEnvironments,
  })
}

export function useRunningEnvironments() {
  return useQuery({
    queryKey: queryKeys.running,
    queryFn: listRunningEnvironments,
    refetchInterval: 10000,
  })
}

// Single environment
export function useEnvironmentInfo(path: string) {
  return useQuery({
    queryKey: queryKeys.info(path),
    queryFn: () => getEnvironmentInfo(path),
    enabled: !!path,
  })
}

export function useEnvironmentStatus(path: string) {
  return useQuery({
    queryKey: queryKeys.status(path),
    queryFn: () => getEnvironmentStatus(path),
    enabled: !!path,
    refetchInterval: 5000,
  })
}

// Mutations
export function useSyncEnvironments() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: syncEnvironments,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.systemStatus })
      queryClient.invalidateQueries({ queryKey: queryKeys.environments })
      toast({
        title: 'Sync completed',
        description: 'Environment list has been updated.',
        variant: 'success',
      })
    },
    onError: (error: Error) => {
      toast({
        title: 'Sync failed',
        description: error.message,
        variant: 'destructive',
      })
    },
  })
}

export function useStartEnvironment() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ path, options }: { path: string; options?: StartRequest }) =>
      startEnvironment(path, options),
    onSuccess: (_, { path }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.downloaded })
      queryClient.invalidateQueries({ queryKey: queryKeys.running })
      queryClient.invalidateQueries({ queryKey: queryKeys.status(path) })
      toast({
        title: 'Environment started',
        description: `${path} is now running.`,
        variant: 'success',
      })
    },
    onError: (error: Error) => {
      toast({
        title: 'Failed to start',
        description: error.message,
        variant: 'destructive',
      })
    },
  })
}

export function useStopEnvironment() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (path: string) => stopEnvironment(path),
    onSuccess: (_, path) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.downloaded })
      queryClient.invalidateQueries({ queryKey: queryKeys.running })
      queryClient.invalidateQueries({ queryKey: queryKeys.status(path) })
      toast({
        title: 'Environment stopped',
        description: `${path} has been stopped.`,
        variant: 'success',
      })
    },
    onError: (error: Error) => {
      toast({
        title: 'Failed to stop',
        description: error.message,
        variant: 'destructive',
      })
    },
  })
}

export function useRestartEnvironment() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (path: string) => restartEnvironment(path),
    onSuccess: (_, path) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.downloaded })
      queryClient.invalidateQueries({ queryKey: queryKeys.running })
      queryClient.invalidateQueries({ queryKey: queryKeys.status(path) })
      toast({
        title: 'Environment restarted',
        description: `${path} has been restarted.`,
        variant: 'success',
      })
    },
    onError: (error: Error) => {
      toast({
        title: 'Failed to restart',
        description: error.message,
        variant: 'destructive',
      })
    },
  })
}

export function useCleanEnvironment() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ path, options }: { path: string; options?: CleanRequest }) =>
      cleanEnvironment(path, options),
    onSuccess: (_, { path }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.downloaded })
      queryClient.invalidateQueries({ queryKey: queryKeys.running })
      queryClient.invalidateQueries({ queryKey: queryKeys.info(path) })
      queryClient.invalidateQueries({ queryKey: queryKeys.status(path) })
      toast({
        title: 'Environment cleaned',
        description: `${path} has been removed.`,
        variant: 'success',
      })
    },
    onError: (error: Error) => {
      toast({
        title: 'Failed to clean',
        description: error.message,
        variant: 'destructive',
      })
    },
  })
}
