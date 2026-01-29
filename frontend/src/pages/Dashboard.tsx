import { Link } from 'react-router-dom'
import {
  RefreshCw,
  List,
  Download,
  Play,
  CheckCircle,
  AlertCircle,
  Loader2,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
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
      icon: List,
      href: '/environments',
    },
    {
      name: 'Downloaded',
      value: downloaded?.total ?? 0,
      icon: Download,
      href: '/downloaded',
    },
    {
      name: 'Running',
      value: running?.total ?? 0,
      icon: Play,
      href: '/running',
    },
  ]

  return (
    <div className="p-6">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground">
          Manage your Vulhub vulnerability environments
        </p>
      </div>

      {/* System Status */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            System Status
            <Button
              variant="outline"
              size="sm"
              onClick={() => syncMutation.mutate()}
              disabled={syncMutation.isPending}
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
        <CardContent>
          {statusLoading ? (
            <p className="text-muted-foreground">Loading...</p>
          ) : status ? (
            <div className="space-y-4">
              <div className="flex items-center gap-2">
                {status.initialized ? (
                  <>
                    <CheckCircle className="h-5 w-5 text-green-500" />
                    <span>Initialized</span>
                  </>
                ) : (
                  <>
                    <AlertCircle className="h-5 w-5 text-yellow-500" />
                    <span>Not initialized</span>
                    <Badge variant="warning">Run vulhub init</Badge>
                  </>
                )}
              </div>
              {status.last_sync_time && (
                <p className="text-sm text-muted-foreground">
                  Last sync: {formatDate(status.last_sync_time)}
                </p>
              )}
              {status.need_sync && (
                <Badge variant="warning">Sync recommended</Badge>
              )}
              <p className="text-xs text-muted-foreground">
                Version: {status.version}
              </p>
            </div>
          ) : (
            <p className="text-muted-foreground">Unable to load status</p>
          )}
        </CardContent>
      </Card>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        {stats.map((stat) => (
          <Link key={stat.name} to={stat.href}>
            <Card className="transition-shadow hover:shadow-md">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  {stat.name}
                </CardTitle>
                <stat.icon className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {envLoading ? '...' : stat.value}
                </div>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>

      {/* Quick Start */}
      <Card className="mt-8">
        <CardHeader>
          <CardTitle>Quick Start</CardTitle>
        </CardHeader>
        <CardContent>
          <ol className="list-inside list-decimal space-y-2 text-sm">
            <li>
              Click <strong>Sync</strong> to update the environment list from
              GitHub
            </li>
            <li>
              Browse <Link to="/environments" className="text-primary underline">All Environments</Link> or search for a
              specific CVE
            </li>
            <li>Click on an environment to view details and start it</li>
            <li>
              View <Link to="/running" className="text-primary underline">Running</Link> environments to manage active
              containers
            </li>
          </ol>
        </CardContent>
      </Card>
    </div>
  )
}
