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
          <Skeleton key={i} className="h-40" />
        ))}
      </div>
    )
  }

  if (environments.length === 0) {
    return (
      <div className="flex h-40 items-center justify-center rounded-lg border border-dashed">
        <p className="text-muted-foreground">No environments found</p>
      </div>
    )
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      {environments.map((env) => (
        <EnvironmentCard key={env.path} environment={env} />
      ))}
    </div>
  )
}
