import { Navigate } from "react-router-dom";


//code dérivé de stackoverflow
export default function PrivateRoute({ user, children }) {
  if (!user) {
    return <Navigate to="/" />;
  }
  if (!user.emailVerified) {
    return <Navigate to="/MailConfirmationPage" />;
  }
  return children;
}