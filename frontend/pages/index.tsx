import { useRouter } from "next/router";
import { useEffect } from "react";

const Home = () => {
  const router = typeof window !== "undefined" ? useRouter() : null;

  useEffect(() => {
    if (router) {
      console.log(router.pathname);
    }
  }, [router]);

  return (
    <div>
      <h1>Home Page</h1>
    </div>
  );
};

export default Home;
