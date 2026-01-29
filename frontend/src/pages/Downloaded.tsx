import { useState, useMemo } from 'react'
import { Link } from 'react-router-dom'
import { Play, Square, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { SearchInput } from '@/components/common/SearchInput'
import { Skeleton } from '@/components/ui/skeleton'
import {
  useDownloadedEnvironments,
  useStartEnvironment,
  useStopEnvironment,
} from '@/hooks/useEnvironments'

export function Downloaded() {
  const [search, setSearch] = useState('')
  const { data, isLoading } = useDownloadedEnvironments()

  const startMutation = useStartEnvironment()
  const stopMutation = useStopEnvironment()

  const filteredEnvironments = useMemo(() => {
    if (!data?.environments) return []

    if (!search) return data.environments

    const searchLower = search.toLowerCase()
    return data.environments.filter(
      (env) =>
        env.environment.path.toLowerCase().includes(searchLower) ||
        env.environment.name.toLowerCase().includes(searchLower)
    )
  }, [data?.environments, search])

  return (
    <div className="p-6">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">Downloaded Environments</h1>
        <p className="text-muted-foreground">
          Manage locally downloaded environments
        </p>
      </div>

      <div className="mb-6 w-full max-w-sm">
        <SearchInput
          value={search}
          onChange={setSearch}
          placeholder="Search downloaded environments..."
        />
      </div>

      <div className="mb-4">
        <p className="text-sm text-muted-foreground">
          {data?.total ?? 0} downloaded environments
        </p>
      </div>

      {isLoading ? (
        <div className="space-y-4">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} className="h-32" />
          ))}
        </div>
      ) : filteredEnvironments.length === 0 ? (
        <div className="flex h-40 items-center justify-center rounded-lg border border-dashed">
          <p className="text-muted-foreground">No downloaded environments</p>
        </div>
      ) : (
        <div className="space-y-4">
          {filteredEnvironments.map((status) => (
            <Card key={status.environment.path}>
              <CardHeader className="pb-2">
                <div className="flex items-start justify-between">
                  <Link
                    to={`/environment/${encodeURIComponent(status.environment.path)}`}
                    className="hover:underline"
                  >
                    <CardTitle className="text-lg">
                      {status.environment.name}
                    </CardTitle>
                  </Link>
                  <Badge variant={status.running ? 'success' : 'secondary'}>
                    {status.running ? 'Running' : 'Stopped'}
                  </Badge>
                </div>
              </CardHeader>
              <CardContent>
                <p className="mb-2 text-sm text-muted-foreground">
                  {status.environment.path}
                </p>
                {status.containers && status.containers.length > 0 && (
                  <div className="mb-3 flex flex-wrap gap-2">
                    {status.containers.map((container) => (
                      <Badge
                        key={container.id}
                        variant="outline"
                        className="text-xs"
                      >
                        {container.name}: {container.state}
                      </Badge>
                    ))}
                  </div>
                )}
                <div className="flex gap-2">
                  {status.running ? (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => stopMutation.mutate(status.environment.path)}
                      disabled={stopMutation.isPending}
                    >
                      {stopMutation.isPending ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <Square className="h-4 w-4" />
                      )}
                      Stop
                    </Button>
                  ) : (
                    <Button
                      size="sm"
                      onClick={() =>
                        startMutation.mutate({ path: status.environment.path })
                      }
                      disabled={startMutation.isPending}
                    >
                      {startMutation.isPending ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <Play className="h-4 w-4" />
                      )}
                      Start
                    </Button>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}
