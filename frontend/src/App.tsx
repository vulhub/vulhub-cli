import { Routes, Route } from 'react-router-dom'
import { Layout } from '@/components/layout/Layout'
import { Dashboard } from '@/pages/Dashboard'
import { Environments } from '@/pages/Environments'
import { Downloaded } from '@/pages/Downloaded'
import { Running } from '@/pages/Running'
import { EnvironmentInfo } from '@/pages/EnvironmentInfo'

function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route path="/" element={<Dashboard />} />
        <Route path="/environments" element={<Environments />} />
        <Route path="/downloaded" element={<Downloaded />} />
        <Route path="/running" element={<Running />} />
        <Route path="/environment/*" element={<EnvironmentInfo />} />
      </Route>
    </Routes>
  )
}

export default App
