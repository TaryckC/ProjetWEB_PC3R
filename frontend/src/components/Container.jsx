import "../styles.css"

function Container({ children, id }) {
    return (
        <div id={id} className="container">
            {children}
        </div>
    )
}

export default Container;