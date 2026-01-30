import { useState, useMemo } from 'react'
import { Play, Download, Package, Bug, Search } from 'lucide-react'
import { SearchInput } from '@/components/common/SearchInput'
import { EnvironmentList } from '@/components/environments/EnvironmentList'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useEnvironments } from '@/hooks/useEnvironments'
import type { FilterType, Environment } from '@/types'

interface GroupedEnvironments {
  running: Environment[]
  downloaded: Environment[]
  available: Environment[]
}

export function Environments() {
  const [search, setSearch] = useState('')
  const [filter, setFilter] = useState<FilterType>('all')

  const { data, isLoading } = useEnvironments()

  const { grouped, totalFiltered } = useMemo(() => {
    if (!data?.environments) {
      return {
        grouped: { running: [], downloaded: [], available: [] } as GroupedEnvironments,
        totalFiltered: 0,
      }
    }

    let filtered = data.environments

    // Apply filter
    if (filter === 'downloaded') {
      filtered = filtered.filter((env) => env.downloaded)
    } else if (filter === 'running') {
      filtered = filtered.filter((env) => env.running)
    }

    // Apply search
    if (search) {
      const searchLower = search.toLowerCase()
      filtered = filtered.filter(
        (env) =>
          env.path.toLowerCase().includes(searchLower) ||
          env.name.toLowerCase().includes(searchLower) ||
          env.app?.toLowerCase().includes(searchLower) ||
          env.cve?.some((cve) => cve.toLowerCase().includes(searchLower)) ||
          env.tags?.some((tag) => tag.toLowerCase().includes(searchLower))
      )
    }

    // Group environments
    const grouped: GroupedEnvironments = {
      running: filtered.filter((env) => env.running).sort((a, b) => a.path.localeCompare(b.path)),
      downloaded: filtered.filter((env) => env.downloaded && !env.running).sort((a, b) => a.path.localeCompare(b.path)),
      available: filtered.filter((env) => !env.downloaded).sort((a, b) => a.path.localeCompare(b.path)),
    }

    return { grouped, totalFiltered: filtered.length }
  }, [data?.environments, filter, search])

  return (
    <div className="relative min-h-full p-6">
      {/* Background grid */}
      <div className="cyber-grid pointer-events-none absolute inset-0 opacity-30" />

      <div className="relative">
        {/* Header */}
        <div className="mb-8">
          <div className="mb-2 flex items-center gap-3">
            <div className="rounded-lg bg-primary/15 p-2">
              <Bug className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h1 className="text-2xl font-bold tracking-tight text-foreground">
                Environments
              </h1>
              <p className="text-sm text-muted-foreground">
                Browse and deploy vulnerability labs
              </p>
            </div>
          </div>
        </div>

        {/* Controls */}
        <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div className="w-full max-w-sm">
            <SearchInput
              value={search}
              onChange={setSearch}
              placeholder="Search CVE, name, app, tag..."
            />
          </div>
          <Tabs
            value={filter}
            onValueChange={(v) => setFilter(v as FilterType)}
          >
            <TabsList className="border border-border/50 bg-card/80">
              <TabsTrigger
                value="all"
                className="text-xs data-[state=active]:bg-primary/15 data-[state=active]:text-primary"
              >
                All
              </TabsTrigger>
              <TabsTrigger
                value="downloaded"
                className="text-xs data-[state=active]:bg-primary/15 data-[state=active]:text-primary"
              >
                Downloaded
              </TabsTrigger>
              <TabsTrigger
                value="running"
                className="text-xs data-[state=active]:bg-primary/15 data-[state=active]:text-primary"
              >
                Running
              </TabsTrigger>
            </TabsList>
          </Tabs>
        </div>

        {/* Stats bar */}
        <div className="mb-6 flex items-center gap-4 rounded-lg border border-border/50 bg-card/50 px-4 py-2 backdrop-blur">
          <div className="flex items-center gap-2">
            <Search className="h-4 w-4 text-muted-foreground" />
            <span className="text-xs text-muted-foreground">Results:</span>
            <span className="text-sm font-medium text-foreground">
              {totalFiltered}
            </span>
          </div>
          <span className="text-border">/</span>
          <div className="flex items-center gap-2">
            <span className="text-xs text-muted-foreground">Total:</span>
            <span className="text-sm text-foreground">
              {data?.total ?? 0}
            </span>
          </div>
        </div>

        {/* Environment Groups */}
        <div className="space-y-8">
          {grouped.running.length > 0 && (
            <section>
              <div className="mb-4 flex items-center gap-3">
                <div className="flex items-center gap-2 rounded-lg bg-emerald-500/15 px-3 py-1.5">
                  <Play className="h-4 w-4 text-emerald-500" />
                  <span className="text-sm font-medium text-emerald-500">
                    Running
                  </span>
                </div>
                <span className="text-xs text-muted-foreground">
                  {grouped.running.length} active
                </span>
                <div className="h-px flex-1 bg-gradient-to-r from-emerald-500/30 to-transparent" />
              </div>
              <EnvironmentList environments={grouped.running} isLoading={isLoading} />
            </section>
          )}

          {grouped.downloaded.length > 0 && (
            <section>
              <div className="mb-4 flex items-center gap-3">
                <div className="flex items-center gap-2 rounded-lg bg-cyan-500/15 px-3 py-1.5">
                  <Download className="h-4 w-4 text-cyan-500" />
                  <span className="text-sm font-medium text-cyan-500">
                    Downloaded
                  </span>
                </div>
                <span className="text-xs text-muted-foreground">
                  {grouped.downloaded.length} ready
                </span>
                <div className="h-px flex-1 bg-gradient-to-r from-cyan-500/30 to-transparent" />
              </div>
              <EnvironmentList environments={grouped.downloaded} isLoading={isLoading} />
            </section>
          )}

          {grouped.available.length > 0 && (
            <section>
              <div className="mb-4 flex items-center gap-3">
                <div className="flex items-center gap-2 rounded-lg bg-muted px-3 py-1.5">
                  <Package className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-medium text-muted-foreground">
                    Available
                  </span>
                </div>
                <span className="text-xs text-muted-foreground">
                  {grouped.available.length} environments
                </span>
                <div className="h-px flex-1 bg-gradient-to-r from-border to-transparent" />
              </div>
              <EnvironmentList environments={grouped.available} isLoading={isLoading} />
            </section>
          )}

          {!isLoading && totalFiltered === 0 && (
            <div className="flex h-40 flex-col items-center justify-center rounded-lg border border-dashed border-border/50 bg-card/30">
              <Search className="mb-2 h-8 w-8 text-muted-foreground" />
              <p className="text-sm text-muted-foreground">
                No environments found
              </p>
              {search && (
                <p className="mt-1 text-xs text-muted-foreground">
                  Try adjusting your search query
                </p>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
