import { NavLink } from 'react-router-dom'
import {
  LayoutDashboard,
  Bug,
  Shield,
  Activity,
} from 'lucide-react'
import { cn } from '@/lib/utils'

const navigation = [
  { name: 'Dashboard', href: '/', icon: LayoutDashboard },
  { name: 'Environments', href: '/environments', icon: Bug },
]

export function Sidebar() {
  return (
    <div className="flex h-full w-64 flex-col border-r border-border/50 bg-card/50 backdrop-blur">
      {/* Logo */}
      <div className="flex h-16 items-center gap-3 border-b border-border/50 px-6">
        <div className="rounded-lg bg-primary/15 p-1.5">
          <Shield className="h-5 w-5 text-primary" />
        </div>
        <span className="text-lg font-semibold tracking-tight text-foreground">
          Vulhub
        </span>
      </div>

      {/* Navigation */}
      <nav className="flex-1 space-y-1 p-4">
        <p className="mb-3 px-3 text-[10px] font-medium uppercase tracking-widest text-muted-foreground">
          Navigation
        </p>
        {navigation.map((item) => (
          <NavLink
            key={item.name}
            to={item.href}
            className={({ isActive }) =>
              cn(
                'group flex cursor-pointer items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary/50',
                isActive
                  ? 'bg-primary/15 text-primary'
                  : 'text-muted-foreground hover:bg-muted hover:text-foreground'
              )
            }
          >
            {({ isActive }) => (
              <>
                <item.icon className="h-4 w-4" />
                <span>{item.name}</span>
                {isActive && (
                  <span className="ml-auto h-1.5 w-1.5 rounded-full bg-primary" />
                )}
              </>
            )}
          </NavLink>
        ))}
      </nav>

      {/* Footer */}
      <div className="border-t border-border/50 p-4">
        <div className="flex items-center gap-2 rounded-lg bg-muted/50 px-3 py-2">
          <Activity className="h-4 w-4 text-emerald-500" />
          <div className="flex-1">
            <p className="text-[10px] font-medium uppercase tracking-wider text-muted-foreground">
              Status
            </p>
            <p className="flex items-center gap-1.5 text-xs text-emerald-500">
              <span className="relative flex h-1.5 w-1.5">
                <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-500 opacity-75" />
                <span className="relative inline-flex h-1.5 w-1.5 rounded-full bg-emerald-500" />
              </span>
              Connected
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
