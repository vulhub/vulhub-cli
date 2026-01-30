import { Box, ExternalLink } from 'lucide-react'
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
      <Card className="border-border/50 bg-card/80 backdrop-blur">
        <CardHeader className="border-b border-border/50">
          <CardTitle className="flex items-center gap-2 font-mono text-sm uppercase tracking-wider">
            <Box className="h-4 w-4 text-primary" />
            Containers
          </CardTitle>
        </CardHeader>
        <CardContent className="pt-4">
          <p className="font-mono text-sm text-muted-foreground">
            No containers running
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="border-border/50 bg-card/80 backdrop-blur">
      <CardHeader className="border-b border-border/50">
        <CardTitle className="flex items-center gap-2 font-mono text-sm uppercase tracking-wider">
          <Box className="h-4 w-4 text-primary" />
          Containers
          <span className="ml-auto font-mono text-xs font-normal text-muted-foreground">
            {containers.length} instance{containers.length !== 1 ? 's' : ''}
          </span>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-3 pt-4">
        {containers.map((container) => (
          <div
            key={container.id}
            className="rounded-lg border border-border/50 bg-muted/30 p-4"
          >
            <div className="flex items-start justify-between">
              <div>
                <p className="font-mono text-sm font-medium text-foreground">
                  {container.name}
                </p>
                <p className="font-mono text-[10px] text-muted-foreground">
                  ID: {container.id.slice(0, 12)}
                </p>
              </div>
              {container.state === 'running' ? (
                <Badge className="border-0 bg-emerald-400/20 font-mono text-[10px] font-normal text-emerald-400">
                  <span className="mr-1.5 inline-block h-1.5 w-1.5 animate-pulse rounded-full bg-emerald-400" />
                  RUNNING
                </Badge>
              ) : (
                <Badge className="border-0 bg-amber-400/20 font-mono text-[10px] font-normal text-amber-400">
                  {container.state?.toUpperCase()}
                </Badge>
              )}
            </div>

            <div className="mt-3 space-y-2 font-mono text-xs">
              <div className="flex items-center gap-2">
                <span className="text-muted-foreground">IMAGE:</span>
                <span className="text-cyan-400">{container.image}</span>
              </div>

              {container.started_at && (
                <div className="flex items-center gap-2">
                  <span className="text-muted-foreground">UPTIME:</span>
                  <span className="text-foreground">
                    {formatUptime(container.started_at)}
                  </span>
                </div>
              )}

              {container.ports && container.ports.length > 0 && (
                <div>
                  <span className="text-muted-foreground">PORTS:</span>
                  <div className="mt-2 flex flex-wrap gap-2">
                    {container.ports.map((port, i) => (
                      <a
                        key={i}
                        href={`http://${port.host_ip || 'localhost'}:${port.host_port}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-1.5 rounded border border-primary/50 bg-primary/10 px-2 py-1 text-[10px] text-primary transition-colors hover:bg-primary/20"
                      >
                        {port.host_port}:{port.container_port}/{port.protocol}
                        <ExternalLink className="h-3 w-3" />
                      </a>
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
