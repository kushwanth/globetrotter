import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import Confetti from "react-confetti";

export default function GeoQuiz() {
  const [question, setQuestion] = useState(null);
  const [selected, setSelected] = useState(null);
  const [score, setScore] = useState(0);
  const [feedback, setFeedback] = useState(null);

  useEffect(() => {
    fetch("https://example.com/api/quiz") // Replace with actual API URL
      .then(res => res.json())
      .then(data => setQuestion(data));
  }, []);

  const handleAnswer = (answer) => {
    setSelected(answer);
    if (answer === question.correctAnswer) {
      setScore(score + 1);
      setFeedback("correct");
    } else {
      setFeedback("incorrect");
    }
  };

  if (!question) return <div>Loading...</div>;

  return (
    <div className="flex flex-col items-center p-6 space-y-6 bg-gray-100 min-h-screen">
      {feedback === "correct" && <Confetti />}
      <div className="bg-white shadow-lg rounded-2xl p-6 w-full max-w-md">
        <h2 className="text-xl font-bold">Guess the City!</h2>
        <p className="mt-2 text-gray-700">{question.clues[Math.floor(Math.random() * question.clues.length)]}</p>
        <div className="mt-4 space-y-2">
          {question.options.map((option, index) => (
            <button
              key={index}
              className={`w-full py-2 px-4 rounded ${selected === option ? "bg-blue-500 text-white" : "bg-gray-200"}`}
              onClick={() => handleAnswer(option)}
            >
              {option}
            </button>
          ))}
        </div>
      </div>

      {feedback && (
        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} className="text-center">
          <h3 className={`text-2xl font-bold ${feedback === "correct" ? "text-green-500" : "text-red-500"}`}>
            {feedback === "correct" ? "🎉 Correct!" : "😢 Incorrect!"}
          </h3>
          <p className="mt-2 text-gray-700">{question.fun_fact}</p>
          <button className="mt-4 bg-blue-500 text-white py-2 px-4 rounded" onClick={() => window.location.reload()}>
            Play Again
          </button>
        </motion.div>
      )}

      <div className="text-lg font-semibold">Score: {score}</div>
    </div>
  );
}
