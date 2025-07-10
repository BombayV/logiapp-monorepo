import Link from "next/link";
import { Button } from "./ui/button";
import Image from "next/image";

export default function Header() {
  return (
    <header className="flex items-center justify-between p-4 border-b bg-background">
      <Link href="/" className="flex items-center space-x-2">
        <Image src="/favicon/favicon-32x32.png" alt="LogiApp Logo" width={32} height={32} />
        <h1 className="text-2xl font-bold">LogiApp</h1>
      </Link>
      <Link href="/login">
        <Button>Iniciar Sesión</Button>
      </Link>
    </header>
  );
}
