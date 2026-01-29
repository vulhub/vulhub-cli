import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { formatUptime } from '@/lib/utils'
import type { ContainerStatus as ContainerStatusType } from '@/types'

interface ContainerStatusProps {
  containers: ContainerStatusType[]
}

export function ContainerStatus({ containers }: ContainerStatusProps) {
  if (!containers || containers.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Containers</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">No containers running</p>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">Containers</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {containers.map((container) => (
          <div
            key={container.id}
            className="rounded-lg border bg-muted/50 p-4"
          >
            <div className="flex items-start justify-between">
              <div>
                <p className="font-medium">{container.name}</p>
                <p className="text-xs text-muted-foreground">
                  {container.id.slice(0, 12)}
                </p>
              </div>
              <Badge
                variant={container.state === 'running' ? 'success' : 'secondary'}
              >
                {container.state}
              </Badge>
            </div>
            <div className="mt-2 space-y-1 text-sm">
              <p>
                <span className="text-muted-foreground">Image:</span>{' '}
                {container.image}
              </p>
              {container.started_at && (
                <p>
                  <span className="text-muted-foreground">Uptime:</span>{' '}
                  {formatUptime(container.started_at)}
                </p>
              )}
              {container.ports && container.ports.length > 0 && (
                <div>
                  <span className="text-muted-foreground">Ports:</span>
                  <div className="mt-1 flex flex-wrap gap-1">
                    {container.ports.map((port, i) => (
                      <Badge key={i} variant="outline" className="text-xs">
                        {port.host_port}:{port.container_port}/{port.protocol}
                      </Badge>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
  )
}
