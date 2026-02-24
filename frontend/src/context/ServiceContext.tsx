import React from "react";
import { createContext } from "react";
import type { IAuthRepository } from "../domain/auth/auth.repository";
import type { IUserRepository } from "../domain/user/user.repository";
import type { IVideoRepository } from "../domain/video/video.repository";
import type { IVideoAnalyzer } from "../domain/video/video.service";
import { AuthRepositoryImpl } from "../gateways/auth/auth.repository.impl";
import { UserRepositoryImpl } from "../gateways/user/user.repository.impl";
import { VideoRepositoryImpl } from "../gateways/video/video.repository.impl";
import { HlsVideoAnalyzer } from "../gateways/video/video.service.impl";
import { VideoRepositoryMock } from "../gateways/video/video.repository.mock";
import { AuthRepositoryMock } from "../gateways/auth/auth.repository.mock";

interface ServiceContextType {
  authRepo: IAuthRepository;
  userRepo: IUserRepository;
  videoRepo: IVideoRepository;
  videoAnalyzer: IVideoAnalyzer;
}

const ServiceContext = createContext<ServiceContextType | undefined>(undefined);

export const ServiceProvider = ({ children }: { children: React.ReactNode }) => {
  const isMock = import.meta.env.VITE_USE_MOCK === "true";

  const services: ServiceContextType = {
    authRepo: isMock ? new AuthRepositoryMock() : new AuthRepositoryImpl(),
    userRepo: new UserRepositoryImpl(),
    videoRepo: isMock ? new VideoRepositoryMock() : new VideoRepositoryImpl(),
    videoAnalyzer: new HlsVideoAnalyzer(),
  };

  return (
    <ServiceContext.Provider value={services}>
      {children}
    </ServiceContext.Provider>
  )
};

// helper hook to use the services in custom hooks or components
export const useServices = (): ServiceContextType => {
  const context = React.useContext(ServiceContext);
  if (!context) {
    throw new Error("useServices must be used within a ServiceProvider");
  }
  return context;
}