"use client"

import { useEffect } from "react"
import { useRouter, usePathname } from "next/navigation"
import Link from "next/link"
import { useStore } from "@/store/use-store"
import { 
  LayoutDashboard, 
  LogOut, 
  User as UserIcon,
  Users
} from "lucide-react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const router = useRouter()
  const pathname = usePathname()
  const user = useStore((state) => state.user)
  const logout = useStore((state) => state.logout)

  useEffect(() => {
    if (!user) {
      router.push('/login')
    } else if (user.role !== 'admin') {
      router.push('/dashboard')
    }
  }, [user, router])

  if (!user || user.role !== 'admin') {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-50">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-slate-200 border-t-slate-900" />
      </div>
    )
  }

  const navigation = [
    { name: 'Заявки', href: '/admin', icon: LayoutDashboard },
    // { name: 'Сотрудники', href: '/admin/users', icon: Users }, // Placeholder
  ]

  return (
    <div className="flex min-h-screen bg-slate-50">
      {/* Sidebar */}
      <div className="fixed inset-y-0 left-0 w-64 bg-slate-900 text-white border-r border-slate-800 p-4 flex flex-col">
        <div className="flex items-center gap-2 px-2 mb-8">
          <div className="h-8 w-8 bg-white rounded-lg flex items-center justify-center text-slate-900 font-bold">
            SB
          </div>
          <span className="text-xl font-bold">Admin Panel</span>
        </div>

        <nav className="flex-1 space-y-1">
          {navigation.map((item) => {
            const isActive = pathname === item.href || pathname.startsWith(item.href + '/')
            return (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  "flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-md transition-colors",
                  isActive 
                    ? "bg-slate-800 text-white" 
                    : "text-slate-400 hover:bg-slate-800 hover:text-white"
                )}
              >
                <item.icon className="h-5 w-5" />
                {item.name}
              </Link>
            )
          })}
        </nav>

        <div className="border-t border-slate-800 pt-4 mt-auto">
          <div className="flex items-center gap-3 px-2 mb-4">
            <div className="h-8 w-8 bg-slate-800 rounded-full flex items-center justify-center">
              <UserIcon className="h-4 w-4 text-slate-400" />
            </div>
            <div className="flex-1 overflow-hidden">
              <p className="text-sm font-medium truncate">{user.name}</p>
              <p className="text-xs text-slate-400 truncate">Администратор</p>
            </div>
          </div>
          <Button 
            variant="ghost" 
            className="w-full justify-start text-red-400 hover:text-red-300 hover:bg-slate-800"
            onClick={() => {
              logout()
              router.push('/login')
            }}
          >
            <LogOut className="h-4 w-4 mr-2" />
            Выйти
          </Button>
        </div>
      </div>

      {/* Main Content */}
      <main className="flex-1 ml-64 p-8">
        <div className="max-w-6xl mx-auto">
          {children}
        </div>
      </main>
    </div>
  )
}

