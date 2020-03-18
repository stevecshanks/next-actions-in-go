import React from "react";
import { render } from "@testing-library/react";
import fetchMock from "jest-fetch-mock";
import App from "./App";

test("renders the actions returned by the API", async () => {
  const response = {
    data: [
      {
        type: "actions",
        id: "12345",
        name: "An action"
      }
    ]
  };
  fetchMock.mockResponse(JSON.stringify(response));

  const { findByText } = render(<App />);

  const action = await findByText("An action");
  expect(action).toBeInTheDocument();
});
