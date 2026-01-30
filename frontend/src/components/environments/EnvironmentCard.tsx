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
      <Card className="group cursor-pointer border-border/50 bg-card/80 backdrop-blur transition-all duration-300 hover:border-primary/50 hover:glow-primary focus-within:border-primary/50 focus-within:glow-primary">
        <CardHeader className="pb-2">
          <div className="flex items-start justify-between gap-2">
            <CardTitle className="font-mono text-sm font-medium tracking-wide text-foreground transition-colors group-hover:text-primary">
              {environment.name}
            </CardTitle>
            <StatusBadge
              running={environment.running}
              downloaded={environment.downloaded}
            />
          </div>
        </CardHeader>
        <CardContent>
          <p className="mb-3 font-mono text-xs text-muted-foreground">
            {environment.path}
          </p>
          <div className="flex flex-wrap gap-1.5">
            {environment.cve?.map((cve) => (
              <Badge
                key={cve}
                variant="destructive"
                className="border-0 bg-red-500/20 font-mono text-[10px] font-normal text-red-400"
              >
                {cve}
              </Badge>
            ))}
            {environment.tags?.slice(0, 3).map((tag) => (
              <Badge
                key={tag}
                variant="outline"
                className="border-border/50 bg-muted/50 font-mono text-[10px] font-normal text-muted-foreground"
              >
                {tag}
              </Badge>
            ))}
          </div>
          {environment.app && (
            <p className="mt-3 flex items-center gap-2 font-mono text-[10px] text-muted-foreground">
              <span className="text-primary">APP:</span>
              <span>{environment.app}</span>
            </p>
          )}
        </CardContent>
      </Card>
    </Link>
  )
}
