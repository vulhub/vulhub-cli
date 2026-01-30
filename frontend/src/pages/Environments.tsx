import { useState, useMemo } from 'react'
import { Play, Download, Package } from 'lucide-react'
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
    <div className="p-6">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">Environments</h1>
        <p className="text-muted-foreground">
          Browse and manage vulnerability environments
        </p>
      </div>

      <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="w-full max-w-sm">
          <SearchInput
            value={search}
            onChange={setSearch}
            placeholder="Search by CVE, name, app, or tag..."
          />
        </div>
        <Tabs
          value={filter}
          onValueChange={(v) => setFilter(v as FilterType)}
        >
          <TabsList>
            <TabsTrigger value="all">All</TabsTrigger>
            <TabsTrigger value="downloaded">Downloaded</TabsTrigger>
            <TabsTrigger value="running">Running</TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      <div className="mb-4">
        <p className="text-sm text-muted-foreground">
          Showing {totalFiltered} of {data?.total ?? 0} environments
        </p>
      </div>

      <div className="space-y-8">
        {grouped.running.length > 0 && (
          <section>
            <div className="mb-4 flex items-center gap-2">
              <Play className="h-5 w-5 text-green-500" />
              <h2 className="text-lg font-semibold">Running</h2>
              <span className="text-sm text-muted-foreground">
                ({grouped.running.length})
              </span>
            </div>
            <EnvironmentList environments={grouped.running} isLoading={isLoading} />
          </section>
        )}

        {grouped.downloaded.length > 0 && (
          <section>
            <div className="mb-4 flex items-center gap-2">
              <Download className="h-5 w-5 text-blue-500" />
              <h2 className="text-lg font-semibold">Downloaded</h2>
              <span className="text-sm text-muted-foreground">
                ({grouped.downloaded.length})
              </span>
            </div>
            <EnvironmentList environments={grouped.downloaded} isLoading={isLoading} />
          </section>
        )}

        {grouped.available.length > 0 && (
          <section>
            <div className="mb-4 flex items-center gap-2">
              <Package className="h-5 w-5 text-muted-foreground" />
              <h2 className="text-lg font-semibold">Available</h2>
              <span className="text-sm text-muted-foreground">
                ({grouped.available.length})
              </span>
            </div>
            <EnvironmentList environments={grouped.available} isLoading={isLoading} />
          </section>
        )}

        {!isLoading && totalFiltered === 0 && (
          <div className="flex h-40 items-center justify-center rounded-lg border border-dashed">
            <p className="text-muted-foreground">No environments found</p>
          </div>
        )}
      </div>
    </div>
  )
}
