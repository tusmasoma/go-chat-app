import React, { ChangeEvent, FormEvent } from "react";
import { Loader2, Send, Plus } from "lucide-react";

type Props = {
  handleInputChange: (e: ChangeEvent<HTMLInputElement> | ChangeEvent<HTMLTextAreaElement>) => void;
  handleSubmit: (e: FormEvent<HTMLFormElement>) => void;
  input: string;
  isLoading: boolean;
  stop: () => void;
};

const InputForm = ({ handleInputChange, handleSubmit, input, isLoading, stop }: Props) => {
  return (
    <form
      onSubmit={(event) => {
        event.preventDefault();
        handleSubmit(event);
      }}
      className="w-full flex flex-row items-center px-4 py-3 border-t border-gray-200 bg-white fixed bottom-0 left-0 right-0"
    >
      <button
        type="button"
        className="rounded-full p-3 bg-gray-200 text-gray-600 mr-2"
      >
        <Plus className="h-5 w-5" />
      </button>

      <input
        type="text"
        placeholder={isLoading ? "Generating . . ." : "Type a message..."}
        value={input}
        disabled={isLoading}
        onChange={handleInputChange}
        className="flex-1 border border-gray-300 rounded-full px-4 py-2 focus:outline-none focus:border-blue-400 text-gray-900"
      />

      <button
        type="submit"
        className="rounded-full p-3 bg-blue-600 text-white ml-2"
      >
        {isLoading ? (
          <Loader2 className="h-5 w-5 animate-spin" />
        ) : (
          <Send className="h-5 w-5" />
        )}
      </button>
    </form>
  );
};

export default InputForm;
