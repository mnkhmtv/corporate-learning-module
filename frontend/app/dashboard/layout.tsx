"use client"

import { useEffect } from "react"
import { useRouter, usePathname } from "next/navigation"
import Link from "next/link"
import { useStore } from "@/store/use-store"
import { 
  LayoutDashboard, 
  PlusCircle, 
  BookOpen, 
  LogOut, 
  User as UserIcon 
} from "lucide-react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const router = useRouter()
  const pathname = usePathname()
  const { user, isAuthLoading } = useStore()
  const logout = useStore((state) => state.logout)

  useEffect(() => {
    if (!isAuthLoading && !user) {
      router.push('/login')
    }
  }, [user, isAuthLoading, router])
  
  if (isAuthLoading || !user) {
    return (
      <div className="flex h-screen items-center justify-center">
        {/* You can replace this with a proper loader/spinner component */}
        <p>Loading...</p>
      </div>
    )
  }

  const navigation = [
    { name: 'Главная', href: '/dashboard', icon: LayoutDashboard },
    { name: 'Новая заявка', href: '/dashboard/requests/create', icon: PlusCircle },
    // { name: 'Мои обучения', href: '/dashboard/learnings', icon: BookOpen }, // Merged into dashboard for simplicity
  ]

  return (
    <div className="flex min-h-screen bg-[#F2F3F7]">
      {/* Sidebar */}
      <div className="fixed inset-y-0 left-0 w-64 bg-white border-r border-slate-200 p-4 flex flex-col">
        <div className="flex items-center gap-2 px-2 mb-8">
          <div
            className="h-8 w-8 rounded-lg flex items-center justify-center text-white font-bold bg-[#9B93FF]"
          >
            SB
          </div>
          <span className="text-xl font-bold text-slate-900 uppercase">SkillBridge</span>
        </div>

        <nav className="flex-1 space-y-1">
          {navigation.map((item) => {
            const isActive = pathname === item.href
            return (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  "flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-xl transition-colors",
                  isActive 
                    ? "bg-slate-100 text-slate-900" 
                    : "text-slate-600 hover:bg-[#F2F3F7] hover:text-slate-900"
                )}
              >
                <item.icon className="h-5 w-5" />
                {item.name}
              </Link>
            )
          })}
        </nav>

        <div className="border-t border-slate-200 pt-4 mt-auto">
          <div className="flex items-center gap-3 px-2 mb-4">
            <div className="h-8 w-8 bg-slate-100 rounded-full flex items-center justify-center">
              <UserIcon className="h-4 w-4 text-slate-500" />
            </div>
            <div className="flex-1 overflow-hidden">
              <p className="text-sm font-medium text-slate-900 truncate">{user.name}</p>
              <p className="text-xs text-slate-500 truncate">{user.role === 'admin' ? 'Администратор' : user.jobTitle}</p>
            </div>
          </div>
          <Button 
            variant="ghost" 
            className="w-full justify-start text-red-500 hover:text-red-600 hover:bg-red-50"
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
        <div className="max-w-5xl mx-auto">
          {children}
        </div>
      </main>
    </div>
  )
}

