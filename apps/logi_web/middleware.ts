import { NextRequest, NextResponse } from 'next/server';

export async function middleware(request: NextRequest) {
  const session = request.cookies.get('session');
  const { pathname } = request.nextUrl;

  // Allow access to login page
  if (pathname.startsWith('/login')) {
    return NextResponse.next();
  }

  // Redirect to login if no session and not on the login page
  if (!session) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  // If there's a session, verify it with the backend
  const response = await fetch('https://poised-lexy-bombay-97849715.koyeb.app/auth/validate', {
    headers: {
      Cookie: `session=${session.value}`,
    },
  });

  // If the session is valid, allow access
  if (response.ok) {
    // If trying to access login page with a valid session, redirect to dashboard
    if (pathname.startsWith('/login')) {
      return NextResponse.redirect(new URL('/dashboard', request.url));
    }
    return NextResponse.next();
  }

  // If the session is invalid, redirect to login
  const loginUrl = new URL('/login', request.url);
  loginUrl.searchParams.set('from', pathname);
  const responseRedirect = NextResponse.redirect(loginUrl);
  responseRedirect.cookies.delete('session'); // Clean up invalid cookie
  return responseRedirect;
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
};
