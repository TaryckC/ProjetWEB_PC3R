import NavBar from "../components/Navbar";
import ChallengePresentation from "../components/ChallengePresentation";

export default function Home() {
  return (
    <div className="flex flex-col">
      <NavBar />
      <div className="flex-1 overflow-hidden flex">
        <ChallengePresentation />
      </div>
    </div>
  );
}
