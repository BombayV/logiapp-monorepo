import { Button } from "@/components/ui/button";
import { ArrowRight, Bot, Truck, Warehouse } from "lucide-react";
import Link from "next/link";

export default function Home() {
  return (
    <div className="flex flex-col min-h-screen bg-background">
      <main className="flex-1">
        <section className="w-full py-12 md:py-24 lg:py-32 xl:py-48 bg-[url('data:image/svg+xml,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20viewBox%3D%220%200%2032%2032%22%20width%3D%2232%22%20height%3D%2232%22%20fill%3D%22none%22%20stroke%3D%22rgb(228%20228%20231%2F.5)%22%3E%3Cpath%20d%3D%22M0%20.5H31.5V32%22%2F%3E%3C%2Fsvg%3E')]">
          <div className="px-4 md:px-6">
            <div className="flex flex-col items-center justify-center space-y-4 text-center">
              <div className="space-y-2">
                <div className="inline-block rounded-lg bg-neutral-200 px-3 py-1 text-sm">
                  Bienvenido a LogiApp
                </div>
                <h1 className="text-3xl font-bold tracking-tighter sm:text-5xl">
                  Tu Plataforma de Logística Inteligente
                </h1>
                <p className="max-w-[900px] text-gray-500 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
                  Maneja pedidos, envíos y almacenes de manera eficiente con nuestra
                  ayuda inteligente.
                </p>
              </div>
              <Link href="/login" className="inline-block">
                <Button size="lg" className="rounded-full px-8">
                  Comienza Ahora
                </Button>
              </Link>
            </div>
          </div>
        </section>
        <section className="w-full py-12 md:py-24 lg:py-32 bg-neutral-50">
          <div className="px-4 md:px-6">
            <div className="flex flex-col items-center justify-center space-y-4 text-center">
              <div className="space-y-2">
                <div className="inline-block rounded-lg bg-neutral-200 px-3 py-1 text-sm">
                  Características Destacadas
                </div>
                <h2 className="text-3xl font-bold tracking-tighter sm:text-4xl">
                  Todo lo que Necesitas para tu Logística
                </h2>
                <p className="max-w-[900px] text-gray-500 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
                  Nuestra plataforma ofrece todas las herramientas necesarias para
                  gestionar tu logística de manera efectiva y eficiente.
                </p>
              </div>
            </div>
            <div className="mx-auto grid max-w-5xl items-center gap-6 py-12 lg:grid-cols-3 lg:gap-12">
              <div className="grid gap-1">
                <div className="flex justify-center">
                  <Truck className="h-12 w-12 text-gray-500" />
                </div>
                <h3 className="text-xl font-bold text-center">
                  Gestión a Tiempo Real
                </h3>
                <p className="text-sm text-gray-500 text-center">
                  Monitorea tus envíos y pedidos en tiempo real para una mejor
                  toma de decisiones.
                </p>
              </div>
              <div className="grid gap-1">
                <div className="flex justify-center">
                  <ArrowRight className="h-12 w-12 text-gray-500" />
                </div>
                <h3 className="text-xl font-bold text-center">
                  Integración Inteligente
                </h3>
                <p className="text-sm text-gray-500 text-center">
                  Conecta con el personal más cercano y optimiza tus rutas de
                  entrega con nuestra tecnologia.
                </p>
              </div>
              <div className="grid gap-1">
                <div className="flex justify-center">
                  <Warehouse className="h-12 w-12 text-gray-500" />
                </div>
                <h3 className="text-xl font-bold text-center">
                  Gestión de Inventario
                </h3>
                <p className="text-sm text-gray-500 text-center">
                  Controla tu inventario de manera eficiente y evita pérdidas con
                  nuestra plataforma.
                </p>
              </div>
            </div>
          </div>
        </section>
      </main>
      <footer className="flex items-center justify-center w-full h-24 border-t">
        <p className="text-gray-500">© 2025 LogiApp. Todos los derechos reservados.</p>
      </footer>
    </div>
  );
}