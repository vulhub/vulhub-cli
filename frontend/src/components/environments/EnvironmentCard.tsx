import { Link } from 'react-router-dom'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { StatusBadge } from '@/components/common/StatusBadge'
import type { Environment } from '@/types'

interface EnvironmentCardProps {
  environment: Environment
}

export function EnvironmentCard({ environment }: EnvironmentCardProps) {
  return (
    <Link to={`/environment/${encodeURIComponent(environment.path)}`}>
      <Card className="transition-shadow hover:shadow-md">
        <CardHeader className="pb-2">
          <div className="flex items-start justify-between">
            <CardTitle className="text-base">{environment.name}</CardTitle>
            <StatusBadge
              running={environment.running}
              downloaded={environment.downloaded}
            />
          </div>
        </CardHeader>
        <CardContent>
          <p className="mb-2 text-sm text-muted-foreground">
            {environment.path}
          </p>
          <div className="flex flex-wrap gap-1">
            {environment.cve?.map((cve) => (
              <Badge key={cve} variant="destructive" className="text-xs">
                {cve}
              </Badge>
            ))}
            {environment.tags?.slice(0, 3).map((tag) => (
              <Badge key={tag} variant="outline" className="text-xs">
                {tag}
              </Badge>
            ))}
          </div>
          {environment.app && (
            <p className="mt-2 text-xs text-muted-foreground">
              App: {environment.app}
            </p>
          )}
        </CardContent>
      </Card>
    </Link>
  )
}
