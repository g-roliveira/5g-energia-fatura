import { SidebarLayout } from "@/components/sidebar-layout"
import { DashboardContent } from "@/components/dashboard-content"

export default function Page() {
  return (
    <SidebarLayout
      breadcrumbs={[
        { label: "Radix Lyra", href: "/" },
        { label: "Dashboard" },
      ]}
    >
      <DashboardContent />
    </SidebarLayout>
  )
}
