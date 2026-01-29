import { useState, useMemo } from 'react'
import { SearchInput } from '@/components/common/SearchInput'
import { EnvironmentList } from '@/components/environments/EnvironmentList'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useEnvironments } from '@/hooks/useEnvironments'
import type { FilterType } from '@/types'

export function Environments() {
  const [search, setSearch] = useState('')
  const [filter, setFilter] = useState<FilterType>('all')

  const { data, isLoading } = useEnvironments()

  const filteredEnvironments = useMemo(() => {
    if (!data?.environments) return []

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

    return filtered
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
          Showing {filteredEnvironments.length} of {data?.total ?? 0}{' '}
          environments
        </p>
      </div>

      <EnvironmentList
        environments={filteredEnvironments}
        isLoading={isLoading}
      />
    </div>
  )
}
