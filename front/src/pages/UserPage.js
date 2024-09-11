import { UserSidebar } from "../components/Sidebar"
import './UserPage.css'

export default function UserPage(props) {
    const handleLogout = () => { window.location.href = "/" }

    return(
        <main className="mh-100">
            <div className="sidebar-root">
                <UserSidebar/>
                <div className="header-root">
                    <div>

                    </div>
                    <div className="page-wrapper">
                        {props.children}
                    </div>
                </div>
            </div>
        </main>
    )
}