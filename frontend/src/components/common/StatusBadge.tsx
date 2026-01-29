import { Badge } from '@/components/ui/badge'

interface StatusBadgeProps {
  running: boolean
  downloaded?: boolean
}

export function StatusBadge({ running, downloaded }: StatusBadgeProps) {
  if (running) {
    return <Badge variant="success">Running</Badge>
  }
  if (downloaded) {
    return <Badge variant="secondary">Downloaded</Badge>
  }
  return <Badge variant="outline">Available</Badge>
}
