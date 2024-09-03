import './Sidebar.css'

export function UserSidebar() {
    return(
        <main>
            <div className="sidebar-logout">
                <SidebarTile icon={"logout"} label={"Logout"} path="/login"/>
            </div>
            <div className="sidebar" style={{padding: '18px 0 14px 0'}}>
                <SidebarTile icon={"home_iot_device"} label={"Devices"} path={"/devices"}/>
                <SidebarTile icon={"electric_meter"} label={"View Energy Usage"} path={"/property/consumption"}/>
                <SidebarTile icon={"share"} label={"Shared"} path={"/user/shared"}/>
                <SidebarDivider/>
                <SidebarTile icon={"add_home"} label={"Add Property"} path={"/property/add"}/>
                <SidebarTile icon={"monitor_weight_gain"} label={"Add Device"} path={"/device/add"}/>
            </div>
        </main>
    )
}

export function SidebarTile({icon, label, path}) {
    const handleClick = () => { window.location.href = path }
    return(
        <div className="sidebar-tile-container" onClick={handleClick}>
            <span className="sidebar-tooltip">{label}</span>
            <span className="material-symbols-outlined icon">{icon}</span>
        </div>
    )
}

export function SidebarDivider() {
    return( <hr style={{width: "calc(100% - 24px)", marginTop: "2px", marginBottom: "2px"}}/> )
}