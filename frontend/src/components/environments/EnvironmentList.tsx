import { EnvironmentCard } from './EnvironmentCard'
import { Skeleton } from '@/components/ui/skeleton'
import type { Environment } from '@/types'

interface EnvironmentListProps {
  environments: Environment[]
  isLoading?: boolean
}

export function EnvironmentList({
  environments,
  isLoading,
}: EnvironmentListProps) {
  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 6 }).map((_, i) => (
          <div
            key={i}
            className="rounded-lg border border-border/50 bg-card/50 p-4"
          >
            <Skeleton className="mb-3 h-5 w-3/4 bg-muted/50" />
            <Skeleton className="mb-4 h-3 w-1/2 bg-muted/50" />
            <div className="flex gap-2">
              <Skeleton className="h-5 w-20 bg-muted/50" />
              <Skeleton className="h-5 w-16 bg-muted/50" />
            </div>
          </div>
        ))}
      </div>
    )
  }

  if (environments.length === 0) {
    return null
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      {environments.map((env) => (
        <EnvironmentCard key={env.path} environment={env} />
      ))}
    </div>
  )
}
