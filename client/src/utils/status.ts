import type { Expense } from '@/types'

export function getStatusClass(status: Expense['status']) {
  switch (status) {
    case 'awaiting_approval':
      return 'bg-yellow-100 text-yellow-800'
    case 'approved':
      return 'bg-green-100 text-green-800'
    case 'rejected':
      return 'bg-red-100 text-red-800'
    case 'completed':
      return 'bg-slate-100 text-gray-800'
    default:
      return 'bg-slate-100 text-gray-800'
  }
}

export function getStatusText(status: Expense['status']) {
  switch (status) {
    case 'awaiting_approval':
      return 'Menunggu Persetujuan'
    case 'approved':
      return 'Disetujui'
    case 'rejected':
      return 'Ditolak'
    case 'completed':
      return 'Selesai'
    default:
      return status
  }
}
