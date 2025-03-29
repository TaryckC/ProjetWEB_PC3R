function Container({ children, id, bgColor = "bg-transparent" }) {
    return (
        <div
            id={id}
            className={`max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 rounded-xl ${bgColor}`}
        >
            {children}
        </div>
    );
}

export default Container;