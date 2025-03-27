import { Navigate } from "react-router-dom";


//code dérivé de stackoverflow
export default function PrivateRoute({ user, children }) {
  console.log("Private route")
  if (!user) {
    console.log("Private route no login")
    return <Navigate to="/" />;
  }
  if(!user.emailVerified){
    console.log("Private route no mail conf")
    return <Navigate to="/MailConfirmationPage" />;
  }
  return children;
}