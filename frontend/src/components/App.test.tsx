import React from "react";
// @ts-ignore See: https://github.com/testing-library/react-testing-library/issues/610
import { render, waitFor } from "@testing-library/react";
import fetchMock from "jest-fetch-mock";
import MockDate from "mockdate";
import App from "./App";

const NOW = new Date(2020, 0, 15, 10, 30, 0);
const ONE_SECOND = 1000;

beforeEach(() => MockDate.set(NOW));

afterEach(() => MockDate.reset());

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

test("correctly handles due by datetimes from the API", async () => {
  const response = {
    data: [
      {
        type: "actions",
        id: "12345",
        name: "An action",
        dueBy: "2020-01-15T10:30:00Z",
      },
      {
        type: "actions",
        id: "23456",
        name: "An action with no due date",
        dueBy: null,
      },
    ],
  };
  fetchMock.mockResponse(JSON.stringify(response));

  const { findByText, queryByText } = render(<App />);

  const action = await findByText("1/15/2020");
  expect(action).toBeInTheDocument();
  const badDueDate = queryByText("1/1/1970");
  expect(badDueDate).not.toBeInTheDocument();
});

test("includes the count of overdue and due soon actions in the window title", async () => {
  const dueSoon = new Date(NOW.getTime() - ONE_SECOND);

  const response = {
    data: [
      {
        type: "actions",
        id: "12345",
        name: "An overdue action",
        dueBy: "2000-01-01T10:30:00Z",
      },
      {
        type: "actions",
        id: "23456",
        name: "A due soon action",
        dueBy: dueSoon.toISOString(),
      },
    ],
  };
  fetchMock.mockResponse(JSON.stringify(response));

  render(<App />);

  await waitFor(() => expect(document.title).toEqual("(2) Next Actions"));
});

test("does not include action count in title if no overdue or due soon actions", async () => {
  const response = {
    data: [
      {
        type: "actions",
        id: "12345",
        name: "An action",
      },
    ],
  };
  fetchMock.mockResponse(JSON.stringify(response));

  render(<App />);

  await waitFor(() => expect(document.title).toEqual("Next Actions"));
});

test("renders an error if the API call fails", async () => {
  fetchMock.mockReject();

  const { findByText } = render(<App />);

  const error = await findByText("An error occurred");
  expect(error).toBeInTheDocument();
});

test("renders errors returned from the API", async () => {
  const response = {
    errors: [
      {
        detail: "a bad thing",
      },
      {
        detail: "another bad thing",
      },
    ],
  };
  fetchMock.mockResponse(JSON.stringify(response), { status: 500 });

  const { findByText } = render(<App />);

  const error = await findByText("An error occurred: a bad thing");
  expect(error).toBeInTheDocument();
  const anotherError = await findByText("An error occurred: another bad thing");
  expect(anotherError).toBeInTheDocument();
});
