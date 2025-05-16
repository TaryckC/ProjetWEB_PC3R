import NavBar from "../components/Navbar";
import ChallengePresentation from "../components/ChallengePresentation";

export default function Home() {
  return (
    <>
      <NavBar />
      <div className="flex flex-row justify-start items-start gap-6 px-4 mt-8">
        <div className="sticky left-0">
          <ChallengePresentation />
        </div>
      </div>
    </>
  );
}
