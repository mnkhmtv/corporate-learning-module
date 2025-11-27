import Link from "next/link"
import { Button } from "@/components/ui/button"

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-b from-slate-50 to-slate-100 p-4 text-center">
      <h1 className="mb-4 text-5xl font-extrabold tracking-tight text-slate-900 uppercase">
        SkillBridge
      </h1>
      <p className="mb-8 max-w-md text-lg text-slate-600">
        Корпоративная платформа для развития навыков и наставничества.
      </p>
      <div className="flex gap-4">
        <Button asChild size="lg">
          <Link href="/login">Войти</Link>
        </Button>
        <Button variant="outline" size="lg" asChild>
          <Link href="/register">Регистрация</Link>
        </Button>
      </div>
    </div>
  )
}
