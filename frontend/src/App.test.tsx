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
        name: "An action",
        dueBy: null,
      },
    ],
  };
  fetchMock.mockResponse(JSON.stringify(response));

  const { findByText } = render(<App />);

  const action = await findByText("An action");
  expect(action).toBeInTheDocument();
});

test("renders an error if the API call fails", async () => {
  fetchMock.mockReject();

  const { findByText } = render(<App />);

  const error = await findByText("An error occurred");
  expect(error).toBeInTheDocument();
});
