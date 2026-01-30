import { Routes, Route } from 'react-router-dom'
import { Layout } from '@/components/layout/Layout'
import { Dashboard } from '@/pages/Dashboard'
import { Environments } from '@/pages/Environments'
import { EnvironmentInfo } from '@/pages/EnvironmentInfo'

function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route path="/" element={<Dashboard />} />
        <Route path="/environments" element={<Environments />} />
        <Route path="/environment/*" element={<EnvironmentInfo />} />
      </Route>
    </Routes>
  )
}

export default App
