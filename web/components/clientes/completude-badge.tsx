import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'

type CompletudeBadgeProps = {
  status: 'complete' | 'partial' | 'failed'
  className?: string
}

const config: Record<CompletudeBadgeProps['status'], { label: string; className: string }> = {
  complete: { label: 'Completo',  className: 'bg-green-100 text-green-800 border-green-200' },
  partial:  { label: 'Parcial',   className: 'bg-yellow-100 text-yellow-800 border-yellow-200' },
  failed:   { label: 'Falhou',    className: 'bg-red-100 text-red-800 border-red-200' },
}

export function CompletudeBadge({ status, className }: CompletudeBadgeProps) {
  const c = config[status]
  return (
    <Badge variant="outline" className={cn(c.className, className)}>
      {c.label}
    </Badge>
  )
}
