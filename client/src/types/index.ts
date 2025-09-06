export interface User {
  id: number
  email: string
  name: string
  role: 'manager' | 'employee'
  created_at: string
}
