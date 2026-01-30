import { Link } from 'react-router-dom'
import {
  RefreshCw,
  Shield,
  Download,
  Play,
  Terminal,
  AlertTriangle,
  Zap,
  Loader2,
  Activity,
  Server,
  Bug,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  useSystemStatus,
  useEnvironments,
  useDownloadedEnvironments,
  useRunningEnvironments,
  useSyncEnvironments,
} from '@/hooks/useEnvironments'
import { formatDate } from '@/lib/utils'

export function Dashboard() {
  const { data: status, isLoading: statusLoading } = useSystemStatus()
  const { data: environments, isLoading: envLoading } = useEnvironments()
  const { data: downloaded } = useDownloadedEnvironments()
  const { data: running } = useRunningEnvironments()

  const syncMutation = useSyncEnvironments()

  const stats = [
    {
      name: 'Total Environments',
      value: environments?.total ?? 0,
      icon: Bug,
      color: 'text-primary',
      bgColor: 'bg-primary/10',
      description: 'Available environments',
    },
    {
      name: 'Downloaded',
      value: downloaded?.total ?? 0,
      icon: Download,
      color: 'text-cyan-500',
      bgColor: 'bg-cyan-500/10',
      description: 'Ready to start',
    },
    {
      name: 'Running',
      value: running?.total ?? 0,
      icon: Activity,
      color: 'text-emerald-500',
      bgColor: 'bg-emerald-500/10',
      description: 'Active instances',
    },
  ]

  return (
    <div className="relative min-h-full p-6">
      {/* Background grid */}
      <div className="cyber-grid pointer-events-none absolute inset-0 opacity-30" />

      <div className="relative">
        {/* Header */}
        <div className="mb-8">
          <div className="mb-2 flex items-center gap-3">
            <div className="rounded-lg bg-primary/15 p-2">
              <Shield className="h-7 w-7 text-primary" />
            </div>
            <div>
              <h1 className="text-2xl font-bold tracking-tight text-foreground">
                Vulhub Dashboard
              </h1>
              <p className="text-sm text-muted-foreground">
                Vulnerability Lab Control Center
              </p>
            </div>
          </div>
        </div>

        {/* System Status */}
        <Card className="border-glow mb-8 border-border/50 bg-card/80 backdrop-blur">
          <CardHeader className="border-b border-border/50 pb-4">
            <CardTitle className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Terminal className="h-5 w-5 text-primary" />
                <span className="text-sm font-medium">System Status</span>
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={() => syncMutation.mutate()}
                disabled={syncMutation.isPending}
                className="text-xs"
              >
                {syncMutation.isPending ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <RefreshCw className="h-4 w-4" />
                )}
                Sync
              </Button>
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-4">
            {statusLoading ? (
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Loader2 className="h-4 w-4 animate-spin" />
                Loading...
              </div>
            ) : status ? (
              <div className="space-y-3 text-sm">
                <div className="flex items-center gap-3">
                  <span className="text-muted-foreground">Status:</span>
                  {status.initialized ? (
                    <span className="flex items-center gap-2 text-emerald-500">
                      <span className="relative flex h-2 w-2">
                        <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-500 opacity-75" />
                        <span className="relative inline-flex h-2 w-2 rounded-full bg-emerald-500" />
                      </span>
                      Online
                    </span>
                  ) : (
                    <span className="flex items-center gap-2 text-amber-500">
                      <AlertTriangle className="h-4 w-4" />
                      Not Initialized
                    </span>
                  )}
                </div>

                {status.last_sync_time && (
                  <div className="flex items-center gap-3">
                    <span className="text-muted-foreground">Last Sync:</span>
                    <span className="text-foreground">
                      {formatDate(status.last_sync_time)}
                    </span>
                  </div>
                )}

                {status.need_sync && (
                  <div className="flex items-center gap-2 text-amber-500">
                    <Zap className="h-4 w-4" />
                    <span>Sync recommended - new data available</span>
                  </div>
                )}

                <div className="flex items-center gap-3">
                  <span className="text-muted-foreground">Version:</span>
                  <span className="text-foreground">{status.version}</span>
                </div>
              </div>
            ) : (
              <div className="flex items-center gap-2 text-sm text-red-500">
                <AlertTriangle className="h-4 w-4" />
                Connection failed
              </div>
            )}
          </CardContent>
        </Card>

        {/* Stats Grid */}
        <div className="mb-8 grid gap-4 md:grid-cols-3">
          {stats.map((stat) => (
            <Link key={stat.name} to="/environments">
              <Card className="stat-card border-border/50 bg-card/80 backdrop-blur">
                <CardContent className="p-6">
                  <div className="flex items-start justify-between">
                    <div>
                      <p className="text-xs font-medium uppercase tracking-wider text-muted-foreground">
                        {stat.name}
                      </p>
                      <p className="mt-2 text-3xl font-bold tracking-tight text-foreground">
                        {envLoading ? (
                          <span className="animate-pulse">--</span>
                        ) : (
                          stat.value
                        )}
                      </p>
                      <p className="mt-1 text-xs text-muted-foreground">
                        {stat.description}
                      </p>
                    </div>
                    <div className={`rounded-lg p-3 ${stat.bgColor}`}>
                      <stat.icon className={`h-5 w-5 ${stat.color}`} />
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>

        {/* Quick Actions */}
        <div className="grid gap-4 md:grid-cols-2">
          {/* Quick Start Guide */}
          <Card className="border-border/50 bg-card/80 backdrop-blur">
            <CardHeader className="border-b border-border/50 pb-4">
              <CardTitle className="flex items-center gap-2">
                <Server className="h-5 w-5 text-primary" />
                <span className="text-sm font-medium">Quick Start</span>
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-4">
              <ol className="space-y-3 text-sm">
                <li className="flex items-start gap-3">
                  <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-primary/15 text-xs font-medium text-primary">
                    1
                  </span>
                  <span className="text-muted-foreground">
                    Click <span className="text-foreground font-medium">Sync</span> to fetch latest environments
                  </span>
                </li>
                <li className="flex items-start gap-3">
                  <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-primary/15 text-xs font-medium text-primary">
                    2
                  </span>
                  <span className="text-muted-foreground">
                    Browse{' '}
                    <Link
                      to="/environments"
                      className="text-primary hover:underline"
                    >
                      Environments
                    </Link>{' '}
                    or search by CVE
                  </span>
                </li>
                <li className="flex items-start gap-3">
                  <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-primary/15 text-xs font-medium text-primary">
                    3
                  </span>
                  <span className="text-muted-foreground">
                    Select target and click <span className="text-emerald-500 font-medium">Start</span>
                  </span>
                </li>
                <li className="flex items-start gap-3">
                  <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-primary/15 text-xs font-medium text-primary">
                    4
                  </span>
                  <span className="text-muted-foreground">
                    Access exposed ports to begin testing
                  </span>
                </li>
              </ol>
            </CardContent>
          </Card>

          {/* Running Environments Preview */}
          <Card className="border-border/50 bg-card/80 backdrop-blur">
            <CardHeader className="border-b border-border/50 pb-4">
              <CardTitle className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <Play className="h-5 w-5 text-emerald-500" />
                  <span className="text-sm font-medium">Active Instances</span>
                </div>
                {(running?.total ?? 0) > 0 && (
                  <span className="flex items-center gap-2 text-xs text-emerald-500">
                    <span className="relative flex h-2 w-2">
                      <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-500 opacity-75" />
                      <span className="relative inline-flex h-2 w-2 rounded-full bg-emerald-500" />
                    </span>
                    {running?.total} running
                  </span>
                )}
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-4">
              {(running?.total ?? 0) === 0 ? (
                <div className="flex flex-col items-center justify-center py-6 text-center">
                  <div className="mb-3 rounded-full bg-muted p-3">
                    <Activity className="h-6 w-6 text-muted-foreground" />
                  </div>
                  <p className="text-sm text-muted-foreground">
                    No active instances
                  </p>
                  <Link
                    to="/environments"
                    className="mt-2 text-xs text-primary hover:underline"
                  >
                    Start an environment →
                  </Link>
                </div>
              ) : (
                <div className="space-y-2">
                  {running?.environments.slice(0, 3).map((env) => (
                    <Link
                      key={env.environment.path}
                      to={`/environment/${encodeURIComponent(env.environment.path)}`}
                      className="flex items-center justify-between rounded-lg bg-muted/50 px-3 py-2 transition-colors hover:bg-muted"
                    >
                      <span className="text-sm text-foreground">
                        {env.environment.name}
                      </span>
                      <span className="flex items-center gap-1 text-xs text-emerald-500">
                        <span className="h-1.5 w-1.5 rounded-full bg-emerald-500" />
                        Running
                      </span>
                    </Link>
                  ))}
                  {(running?.total ?? 0) > 3 && (
                    <Link
                      to="/environments"
                      className="block text-center text-xs text-muted-foreground hover:text-primary"
                    >
                      +{(running?.total ?? 0) - 3} more →
                    </Link>
                  )}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
