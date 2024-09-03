import "./SmallCard.css"

export default function SmallCard({title, number, measurement}) {

    return(
        <div className="card">
            <p className="smallcard-title">{title}</p>
            <div >
                <p className="smallcard-number">{number}</p>
                <p className="smallcard-measurement">{measurement}</p>
            </div> 
        </div>
    )
}