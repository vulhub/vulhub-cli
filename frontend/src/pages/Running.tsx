import { Link } from 'react-router-dom'
import { Square, RotateCcw, Loader2, ExternalLink } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import {
  useRunningEnvironments,
  useStopEnvironment,
  useRestartEnvironment,
} from '@/hooks/useEnvironments'
import { formatUptime } from '@/lib/utils'

export function Running() {
  const { data, isLoading } = useRunningEnvironments()

  const stopMutation = useStopEnvironment()
  const restartMutation = useRestartEnvironment()

  return (
    <div className="p-6">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">Running Environments</h1>
        <p className="text-muted-foreground">
          Monitor and control active environments
        </p>
      </div>

      <div className="mb-4">
        <p className="text-sm text-muted-foreground">
          {data?.total ?? 0} running environments
        </p>
      </div>

      {isLoading ? (
        <div className="space-y-4">
          {Array.from({ length: 2 }).map((_, i) => (
            <Skeleton key={i} className="h-48" />
          ))}
        </div>
      ) : data?.environments.length === 0 ? (
        <div className="flex h-40 items-center justify-center rounded-lg border border-dashed">
          <div className="text-center">
            <p className="text-muted-foreground">No running environments</p>
            <Link
              to="/environments"
              className="text-sm text-primary hover:underline"
            >
              Start an environment
            </Link>
          </div>
        </div>
      ) : (
        <div className="space-y-4">
          {data?.environments.map((status) => (
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
                  <Badge variant="success">Running</Badge>
                </div>
                <p className="text-sm text-muted-foreground">
                  {status.environment.path}
                </p>
              </CardHeader>
              <CardContent>
                {/* Containers */}
                <div className="mb-4 space-y-2">
                  {status.containers?.map((container) => (
                    <div
                      key={container.id}
                      className="rounded-lg bg-muted/50 p-3"
                    >
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="font-medium">{container.name}</p>
                          <p className="text-xs text-muted-foreground">
                            {container.image}
                          </p>
                        </div>
                        {container.started_at && (
                          <p className="text-sm text-muted-foreground">
                            Uptime: {formatUptime(container.started_at)}
                          </p>
                        )}
                      </div>
                      {container.ports && container.ports.length > 0 && (
                        <div className="mt-2 flex flex-wrap gap-2">
                          {container.ports.map((port, i) => (
                            <a
                              key={i}
                              href={`http://${port.host_ip || 'localhost'}:${port.host_port}`}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="inline-flex items-center gap-1"
                            >
                              <Badge
                                variant="outline"
                                className="cursor-pointer text-xs hover:bg-accent"
                              >
                                {port.host_port}:{port.container_port}
                                <ExternalLink className="ml-1 h-3 w-3" />
                              </Badge>
                            </a>
                          ))}
                        </div>
                      )}
                    </div>
                  ))}
                </div>

                {/* Actions */}
                <div className="flex gap-2">
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
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() =>
                      restartMutation.mutate(status.environment.path)
                    }
                    disabled={restartMutation.isPending}
                  >
                    {restartMutation.isPending ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <RotateCcw className="h-4 w-4" />
                    )}
                    Restart
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}
