import { Navigate, Outlet } from "react-router-dom";
import { useSelector } from "react-redux";

const ProtectedRoutes = ({ allowRoles }) => {
    const { isAuthenticated, role } = useSelector((state) => state.authReducer);

    if (!isAuthenticated) {
        return <Navigate to="/login" replace />
    }

    if (allowRoles && !allowRoles.includes(role)) {
        return <Navigate to="/unauthorized" replace />
    }

    return <Outlet />;
}

export default ProtectedRoutes;