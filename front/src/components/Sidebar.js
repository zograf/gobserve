import './Sidebar.css'

export function UserSidebar({ref1, ref2, ref3, ref4}) {
    return(
        <main>
            <div className="sidebar-logout">
                <SidebarTile icon={"logout"} label={"Logout"} path="/login"/>
            </div>
            <div className="sidebar" style={{padding: '18px 0 14px 0'}}>
                <SidebarTile icon={"view_list"} label={"Services"} reff={ref4}/>
                <SidebarTile icon={"assessment"} label={"Analytics"} reff={ref1}/>
                <SidebarTile icon={"show_chart"} label={"Charts"} reff={ref2}/>
                <SidebarTile icon={"table_chart"} label={"Tables"} reff={ref3}/>
                <SidebarDivider/>
            </div>
        </main>
    )
}

export function SidebarTile({icon, label, reff}) {
    const handleClick = () => {reff.current.scrollIntoView({behavior: 'smooth'})}
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