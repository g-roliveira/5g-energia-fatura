import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'

type StatusBadgeProps = {
  status: 'ativo' | 'inativo' | 'prospecto' | 'succeeded' | 'partial' | 'failed' | 'running'
  className?: string
}

const statusConfig: Record<StatusBadgeProps['status'], { label: string; className: string }> = {
  ativo:      { label: 'Ativo',         className: 'bg-green-100 text-green-800 border-green-200' },
  succeeded:  { label: 'Concluído',     className: 'bg-green-100 text-green-800 border-green-200' },
  inativo:    { label: 'Inativo',       className: 'bg-red-100 text-red-800 border-red-200' },
  failed:     { label: 'Falhou',        className: 'bg-red-100 text-red-800 border-red-200' },
  prospecto:  { label: 'Prospecto',     className: 'bg-yellow-100 text-yellow-800 border-yellow-200' },
  partial:    { label: 'Parcial',       className: 'bg-yellow-100 text-yellow-800 border-yellow-200' },
  running:    { label: 'Sincronizando', className: 'bg-blue-100 text-blue-800 border-blue-200' },
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const config = statusConfig[status] ?? { label: status, className: 'bg-gray-100 text-gray-700' }
  return (
    <Badge variant="outline" className={cn(config.className, className)}>
      {config.label}
    </Badge>
  )
}
