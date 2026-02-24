import { useAuth } from "@/context/AuthContext";
import { useEffect } from "react";
import { useNavigate } from "react-router";

export const LoginPage = () => {
  const { session, login } = useAuth();
  const navigate = useNavigate();
  // 現在はDummyなので、ログイン処理をしてすぐにリダイレクトする
  useEffect(() => {
    if (session.status === "authenticated") {
      // 既にユーザー情報があるなら、トップへ飛ばして終了
      navigate("/");
      return;
    }

    // まだならログインを実行
    console.log("Logging in...");
    login("", "");
  }, [session, login, navigate]);

  return <div>Logging in...</div>;
};