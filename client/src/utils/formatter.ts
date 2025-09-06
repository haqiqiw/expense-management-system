import { format } from 'date-fns'

export function formatRupiah(amount: number): string {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(amount)
}

export function formatDate(dateString: string): string {
  try {
    return format(new Date(dateString), 'd/M/yyyy HH:mm')
  } catch (error) {
    console.log('Failed to format date:', error)
    return 'Invalid Date: '
  }
}
